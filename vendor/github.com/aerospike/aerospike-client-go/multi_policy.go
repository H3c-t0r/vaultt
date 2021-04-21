// Copyright 2013-2020 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import "time"

// MultiPolicy contains parameters for policy attributes used in
// query and scan operations.
type MultiPolicy struct {
	BasePolicy

	// Maximum number of concurrent requests to server nodes at any point in time.
	// If there are 16 nodes in the cluster and maxConcurrentNodes is 8, then queries
	// will be made to 8 nodes in parallel.  When a query completes, a new query will
	// be issued until all 16 nodes have been queried.
	// Default (0) is to issue requests to all server nodes in parallel.
	MaxConcurrentNodes int

	// FailOnClusterChange determines scan termination if cluster is in fluctuating state.
	FailOnClusterChange bool

	// RecordsPerSecond limits returned records per second (rps) rate for each server.
	// Will not apply rps limit if recordsPerSecond is zero (default).
	// Currently only applicable to a query without a defined filter.
	RecordsPerSecond int

	// Number of records to place in queue before blocking.
	// Records received from multiple server nodes will be placed in a queue.
	// A separate goroutine consumes these records in parallel.
	// If the queue is full, the producer goroutines will block until records are consumed.
	RecordQueueSize int //= 50

	// Indicates if bin data is retrieved. If false, only record digests are retrieved.
	IncludeBinData bool //= true;
}

// NewMultiPolicy initializes a MultiPolicy instance with default values.
func NewMultiPolicy() *MultiPolicy {
	bp := *NewPolicy()
	bp.SocketTimeout = 30 * time.Second

	return &MultiPolicy{
		BasePolicy:          bp,
		MaxConcurrentNodes:  0,
		RecordsPerSecond:    0,
		RecordQueueSize:     50,
		IncludeBinData:      true,
		FailOnClusterChange: false,
	}
}
