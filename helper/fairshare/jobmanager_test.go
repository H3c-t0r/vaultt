package fairshare

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestJobManager_NewJobManager(t *testing.T) {
	testCases := []struct {
		name               string
		numWorkers         int
		expectedNumWorkers int
	}{
		{
			name:               "",
			numWorkers:         0,
			expectedNumWorkers: 1,
		},
		{
			name:               "",
			numWorkers:         5,
			expectedNumWorkers: 5,
		},
		{
			name:               "",
			numWorkers:         5,
			expectedNumWorkers: 5,
		},
		{
			name:               "",
			numWorkers:         5,
			expectedNumWorkers: 5,
		},
		{
			name:               "",
			numWorkers:         5,
			expectedNumWorkers: 5,
		},
	}

	l := newTestLogger("jobmanager-test")
	for tcNum, tc := range testCases {
		j := NewJobManager(tc.name, tc.numWorkers, l, nil)

		if tc.name != "" && tc.name != j.name {
			t.Errorf("tc %d: expected name %s, got %s", tcNum, tc.name, j.name)
		}
		if j.queues == nil {
			t.Errorf("tc %d: queues not set up properly", tcNum)
		}
		if j.quit == nil {
			t.Errorf("tc %d: quit channel not set up properly", tcNum)
		}
		if j.workerPool.numWorkers != tc.expectedNumWorkers {
			t.Errorf("tc %d: expected %d workers, got %d", tcNum, tc.expectedNumWorkers, j.workerPool.numWorkers)
		}
		if j.logger == nil {
			t.Errorf("tc %d: logger not set up properly", tcNum)
		}
	}
}

func TestJobManager_Start(t *testing.T) {
	numJobs := 10
	j := NewJobManager("job-mgr-test", 3, newTestLogger("jobmanager-test"), nil)

	var wg sync.WaitGroup
	wg.Add(numJobs)
	j.Start()
	defer j.Stop()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()

	ex := func(_ string) error {
		wg.Done()
		return nil
	}
	onFail := func(_ error) {}

	for i := 0; i < numJobs; i++ {
		// distribute jobs between 3 queues in the job manager
		job := newTestJob(t, fmt.Sprintf("test-job-%d", i), ex, onFail)
		j.AddJob(&job, fmt.Sprintf("queue-%d", i%3))
	}

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}
}

func TestJobManager_StartAndPause(t *testing.T) {
	numJobs := 10
	j := NewJobManager("job-mgr-test", 3, newTestLogger("jobmanager-test"), nil)

	var wg sync.WaitGroup
	wg.Add(numJobs)
	j.Start()
	defer j.Stop()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()

	ex := func(_ string) error {
		wg.Done()
		return nil
	}
	onFail := func(_ error) {}

	for i := 0; i < numJobs; i++ {
		// distribute jobs between 3 queues in the job manager
		job := newTestJob(t, fmt.Sprintf("test-job-%d", i), ex, onFail)
		j.AddJob(&job, fmt.Sprintf("queue-%d", i%3))
	}

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}

	// now that the work queue is empty, let's add more jobs and make sure
	// we pick up where we left off

	for i := 0; i < 5; i++ {
		numAdditionalJobs := 5
		wg.Add(numAdditionalJobs)

		timeout = time.After(5 * time.Second)
		go func() {
			wg.Wait()
			doneCh <- struct{}{}
		}()

		for i := numJobs; i < numJobs+numAdditionalJobs; i++ {
			// distribute jobs between 3 queues in the job manager
			job := newTestJob(t, fmt.Sprintf("test-job-%d", i), ex, onFail)
			j.AddJob(&job, fmt.Sprintf("queue-%d", i%3))
		}

		select {
		case <-doneCh:
			break
		case <-timeout:
			t.Fatal("timed out")
		}

		numJobs += numAdditionalJobs
	}
}

func TestJobManager_Stop(t *testing.T) {
	j := NewJobManager("job-mgr-test", 5, newTestLogger("jobmanager-test"), nil)

	j.Start()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)
	go func() {
		j.Stop()
		j.wg.Wait()
		doneCh <- struct{}{}
	}()

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}
}

func TestFairshare_StopMultiple(t *testing.T) {
	j := NewJobManager("job-mgr-test", 5, newTestLogger("jobmanager-test"), nil)

	j.Start()

	doneCh := make(chan struct{})
	timeout := time.After(5 * time.Second)
	go func() {
		j.Stop()
		j.wg.Wait()
		doneCh <- struct{}{}
	}()

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}

	// essentially, we don't want to panic here
	var r interface{}
	go func() {
		t.Helper()

		defer func() {
			r = recover()
			doneCh <- struct{}{}
		}()

		j.Stop()
		j.wg.Wait()
	}()

	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}

	if r != nil {
		t.Fatalf("panic during second stop: %v", r)
	}
}

func TestJobManager_AddJob(t *testing.T) {
	testCases := []struct {
		name    string
		queueID string
	}{
		{
			name:    "test1",
			queueID: "q1",
		},
		{
			name:    "test2",
			queueID: "q1",
		},
		{
			name:    "test3",
			queueID: "q1",
		},
		{
			name:    "test4",
			queueID: "q2",
		},
		{
			name:    "test5",
			queueID: "q3",
		},
	}

	j := NewJobManager("job-mgr-test", 3, newTestLogger("jobmanager-test"), nil)

	expectedCount := make(map[string]int)
	for _, tc := range testCases {
		if _, ok := expectedCount[tc.queueID]; !ok {
			expectedCount[tc.queueID] = 1
		} else {
			expectedCount[tc.queueID]++
		}

		job := newDefaultTestJob(t, tc.name)
		j.AddJob(&job, tc.queueID)
	}

	if len(expectedCount) != len(j.queues) {
		t.Fatalf("expected %d queues, got %d", len(expectedCount), len(j.queues))
	}

	for k, v := range j.queues {
		if v.Len() != expectedCount[k] {
			t.Fatalf("queue %s has bad count. expected %d, got %d", k, expectedCount[k], v.Len())
		}
	}
}

func TestJobManager_GetPendingJobCount(t *testing.T) {
	numJobs := 15
	j := NewJobManager("test-job-mgr", 3, newTestLogger("jobmanager-test"), nil)

	for i := 0; i < numJobs; i++ {
		job := newDefaultTestJob(t, fmt.Sprintf("job-%d", i))
		j.AddJob(&job, fmt.Sprintf("queue-%d", i%4))
	}

	pendingJobs := j.GetPendingJobCount()
	if pendingJobs != numJobs {
		t.Errorf("expected %d jobs, got %d", numJobs, pendingJobs)
	}
}

func TestJobManager_GetWorkQueueLengths(t *testing.T) {
	j := NewJobManager("test-job-mgr", 3, newTestLogger("jobmanager-test"), nil)

	expected := make(map[string]int)
	for i := 0; i < 25; i++ {
		queueID := fmt.Sprintf("queue-%d", i%4)
		job := newDefaultTestJob(t, fmt.Sprintf("job-%d", i))

		j.AddJob(&job, queueID)

		if _, ok := expected[queueID]; !ok {
			expected[queueID] = 0
		}

		expected[queueID]++
	}

	pendingJobs := j.GetWorkQueueLengths()
	if !reflect.DeepEqual(pendingJobs, expected) {
		t.Errorf("expected %v job count, got %v", expected, pendingJobs)
	}
}

func TestJobManager_EndToEnd(t *testing.T) {
	testCases := []struct {
		name    string
		queueID string
	}{
		{
			name:    "job-1",
			queueID: "queue-1",
		},
		{
			name:    "job-2",
			queueID: "queue-2",
		},
		{
			name:    "job-3",
			queueID: "queue-1",
		},
		{
			name:    "job-4",
			queueID: "queue-3",
		},
		{
			name:    "job-5",
			queueID: "queue-3",
		},
	}

	j := NewJobManager("test-job-mgr", 1, newTestLogger("jobmanager-test"), nil)

	var wg sync.WaitGroup
	ex := func(name string) error {
		wg.Done()
		return nil
	}
	onFail := func(_ error) {}

	for _, tc := range testCases {
		wg.Add(1)
		job := newTestJob(t, tc.name, ex, onFail)
		j.AddJob(&job, tc.queueID)
	}

	j.Start()
	defer j.Stop()

	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()

	timeout := time.After(5 * time.Second)
	select {
	case <-doneCh:
		break
	case <-timeout:
		t.Fatal("timed out")
	}
}

func TestFairshare_nilLoggerJobManager(t *testing.T) {
	j := NewJobManager("test-job-mgr", 1, nil, nil)
	if j.logger == nil {
		t.Error("logger not set up properly")
	}
}

func TestFairshare_getNextQueue(t *testing.T) {
	j := NewJobManager("test-job-mgr", 18, nil, nil)

	for i := 0; i < 10; i++ {
		job := newDefaultTestJob(t, fmt.Sprintf("job-%d", i))
		j.AddJob(&job, "a")
		j.AddJob(&job, "b")
		j.AddJob(&job, "c")
	}

	// fake out some number of workers with various remaining work scenario
	j.l.Lock()
	j.workerCount["a"] = 1
	j.workerCount["b"] = 2
	j.workerCount["c"] = 5
	j.l.Unlock()

	// this also tests j.sortByNumWorkers() indirectly
	expectedOrder := []string{"a", "a", "b", "a", "b", "a", "b", "a", "b", "c"}

	for _, expectedQueueID := range expectedOrder {
		queueID, canAssignWorker := j.getNextQueue()

		if !canAssignWorker {
			t.Fatalf("expected have work true, got false for queue %q", queueID)
		}
		if queueID != expectedQueueID {
			t.Errorf("expected queueID %q, got %q", expectedQueueID, queueID)
		}

		// simulate a worker being added to that queue
		j.workerCount[queueID]++
	}
}

func TestJobManager_pruneEmptyQueues(t *testing.T) {
	j := NewJobManager("test-job-mgr", 18, nil, nil)

	// add a few jobs to test out queue pruning
	// for test simplicity, we'll keep the number of workers per queue at 0
	testJob := newDefaultTestJob(t, "job-0")
	j.AddJob(&testJob, "a")
	j.AddJob(&testJob, "a")
	j.AddJob(&testJob, "b")

	job, queueID := j.getNextJob()
	if queueID != "a" || job == nil {
		t.Fatalf("bad next job: queueID %s, job: %#v", queueID, job)
	}

	j.l.RLock()
	if _, ok := j.queues["a"]; !ok {
		t.Error("expected queue 'a' to exist")
	}
	if _, ok := j.queues["b"]; !ok {
		t.Error("expected queue 'b' to exist")
	}
	j.l.RUnlock()

	job, queueID = j.getNextJob()
	if queueID != "a" || job == nil {
		t.Fatalf("bad next job: queueID %s, job: %#v", queueID, job)
	}

	j.l.RLock()
	if _, ok := j.queues["a"]; ok {
		t.Error("expected queue 'a' to be pruned")
	}
	if _, ok := j.queues["b"]; !ok {
		t.Error("expected queue 'b' to exist")
	}
	j.l.RUnlock()

	job, queueID = j.getNextJob()
	if queueID != "b" || job == nil {
		t.Fatalf("bad next job: queueID %s, job: %#v", queueID, job)
	}

	j.l.RLock()
	if _, ok := j.queues["a"]; ok {
		t.Error("expected queue 'a' to be pruned")
	}
	if _, ok := j.queues["b"]; ok {
		t.Error("expected queue 'b' to be pruned")
	}
	j.l.RUnlock()

	job, queueID = j.getNextJob()
	if job != nil {
		t.Errorf("expected no more jobs (out of queues). queueID: %s, job: %#v", queueID, job)
	}
}

func TestFairshare_WorkerCount_IncrementAndDecrement(t *testing.T) {
	j := NewJobManager("test-job-mgr", 18, nil, nil)

	job := newDefaultTestJob(t, "job-0")
	j.AddJob(&job, "a")
	j.AddJob(&job, "b")
	j.AddJob(&job, "c")

	// test to make sure increment works
	j.incrementWorkerCount("a")
	workerCounts := j.GetWorkerCounts()
	if workerCounts["a"] != 1 {
		t.Fatalf("expected 1 worker on 'a', got %d", workerCounts["a"])
	}
	if workerCounts["b"] != 0 {
		t.Fatalf("expected 0 workers on 'b', got %d", workerCounts["b"])
	}
	if workerCounts["c"] != 0 {
		t.Fatalf("expected 0 workers on 'c', got %d", workerCounts["c"])
	}

	// test to make sure decrement works (when there is still work for the queue)
	j.decrementWorkerCount("a")
	workerCounts = j.GetWorkerCounts()
	if workerCounts["a"] != 0 {
		t.Fatalf("expected 0 workers on 'a', got %d", workerCounts["a"])
	}

	// add a worker to queue "a" and remove all work to ensure worker count gets
	// cleared out for "a"
	j.incrementWorkerCount("a")
	j.l.Lock()
	j.queues["a"].Remove(j.queues["a"].Front())
	j.l.Unlock()

	j.decrementWorkerCount("a")
	workerCounts = j.GetWorkerCounts()
	if _, ok := workerCounts["a"]; ok {
		t.Fatalf("expected no worker count for 'a', got %#v", workerCounts)
	}
}

func TestFairshare_queueWorkersSaturated(t *testing.T) {
	j := NewJobManager("test-job-mgr", 20, nil, nil)

	job := newDefaultTestJob(t, "job-0")
	j.AddJob(&job, "a")
	j.AddJob(&job, "b")

	// no more than 9 workers can be assigned to a single queue in this example
	for i := 0; i < 8; i++ {
		j.incrementWorkerCount("a")
		j.incrementWorkerCount("b")

		j.l.RLock()
		if j.queueWorkersSaturated("a") {
			j.l.RUnlock()
			t.Fatalf("queue 'a' falsely saturated: %#v", j.GetWorkerCounts())
		}
		if j.queueWorkersSaturated("b") {
			j.l.RUnlock()
			t.Fatalf("queue 'b' falsely saturated: %#v", j.GetWorkerCounts())
		}
		j.l.RUnlock()
	}

	// adding the 9th and 10th workers should saturate the number of workers we
	// can have per queue
	for i := 8; i < 10; i++ {
		j.incrementWorkerCount("a")
		j.incrementWorkerCount("b")

		j.l.RLock()
		if !j.queueWorkersSaturated("a") {
			j.l.RUnlock()
			t.Fatalf("queue 'a' falsely unsaturated: %#v", j.GetWorkerCounts())
		}
		if !j.queueWorkersSaturated("b") {
			j.l.RUnlock()
			t.Fatalf("queue 'b' falsely unsaturated: %#v", j.GetWorkerCounts())
		}
		j.l.RUnlock()
	}
}
