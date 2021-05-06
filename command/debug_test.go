package command

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mholt/archiver"
	"github.com/mitchellh/cli"
)

func testDebugCommand(tb testing.TB) (*cli.MockUi, *DebugCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &DebugCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestDebugCommand_Run(t *testing.T) {
	t.Parallel()

	testDir, err := ioutil.TempDir("", "vault-debug")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"valid",
			[]string{
				"-duration=1s",
				fmt.Sprintf("-output=%s/valid", testDir),
			},
			"",
			0,
		},
		{
			"too_many_args",
			[]string{
				"-duration=1s",
				fmt.Sprintf("-output=%s/too_many_args", testDir),
				"foo",
			},
			"Too many arguments",
			1,
		},
		{
			"invalid_target",
			[]string{
				"-duration=1s",
				fmt.Sprintf("-output=%s/invalid_target", testDir),
				"-target=foo",
			},
			"Ignoring invalid targets: foo",
			0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testDebugCommand(t)
			cmd.client = client
			cmd.skipTimingChecks = true

			code := cmd.Run(tc.args)
			if code != tc.code {
				t.Errorf("expected %d to be %d", code, tc.code)
			}

			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, tc.out) {
				t.Fatalf("expected %q to contain %q", combined, tc.out)
			}
		})
	}
}

func TestDebugCommand_Archive(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		ext         string
		expectError bool
	}{
		{
			"no-ext",
			"",
			false,
		},
		{
			"with-ext-tar-gz",
			".tar.gz",
			false,
		},
		{
			"with-ext-tgz",
			".tgz",
			false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create temp dirs for each test case since os.Stat and tgz.Walk
			// (called down below) exhibits raciness otherwise.
			testDir, err := ioutil.TempDir("", "vault-debug")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(testDir)

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testDebugCommand(t)
			cmd.client = client
			cmd.skipTimingChecks = true

			// We use tc.name as the base path and apply the extension per
			// test case.
			basePath := tc.name
			outputPath := filepath.Join(testDir, basePath+tc.ext)
			args := []string{
				"-duration=1s",
				fmt.Sprintf("-output=%s", outputPath),
				"-target=server-status",
			}

			code := cmd.Run(args)
			if exp := 0; code != exp {
				t.Log(ui.OutputWriter.String())
				t.Log(ui.ErrorWriter.String())
				t.Fatalf("expected %d to be %d", code, exp)
			}
			// If we expect an error we're done here
			if tc.expectError {
				return
			}

			expectedExt := tc.ext
			if expectedExt == "" {
				expectedExt = debugCompressionExt
			}

			bundlePath := filepath.Join(testDir, basePath+expectedExt)
			_, err = os.Stat(bundlePath)
			if os.IsNotExist(err) {
				t.Log(ui.OutputWriter.String())
				t.Fatal(err)
			}

			tgz := archiver.NewTarGz()
			err = tgz.Walk(bundlePath, func(f archiver.File) error {
				fh, ok := f.Header.(*tar.Header)
				if !ok {
					t.Fatalf("invalid file header: %#v", f.Header)
				}

				// Ignore base directory and index file
				if fh.Name == basePath+"/" || fh.Name == filepath.Join(basePath, "index.json") {
					return nil
				}

				if fh.Name != filepath.Join(basePath, "server_status.json") {
					t.Fatalf("unxexpected file: %s", fh.Name)
				}
				return nil
			})
		})
	}
}

func TestDebugCommand_CaptureTargets(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		targets       []string
		expectedFiles []string
	}{
		{
			"config",
			[]string{"config"},
			[]string{"config.json"},
		},
		{
			"host-info",
			[]string{"host"},
			[]string{"host_info.json"},
		},
		{
			"metrics",
			[]string{"metrics"},
			[]string{"metrics.json"},
		},
		{
			"replication-status",
			[]string{"replication-status"},
			[]string{"replication_status.json"},
		},
		{
			"server-status",
			[]string{"server-status"},
			[]string{"server_status.json"},
		},
		{
			"all-minus-pprof",
			[]string{"config", "host", "metrics", "replication-status", "server-status"},
			[]string{"config.json", "host_info.json", "metrics.json", "replication_status.json", "server_status.json"},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testDir, err := ioutil.TempDir("", "vault-debug")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(testDir)

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testDebugCommand(t)
			cmd.client = client
			cmd.skipTimingChecks = true

			basePath := tc.name
			args := []string{
				"-duration=1s",
				fmt.Sprintf("-output=%s/%s", testDir, basePath),
			}
			for _, target := range tc.targets {
				args = append(args, fmt.Sprintf("-target=%s", target))
			}

			code := cmd.Run(args)
			if exp := 0; code != exp {
				t.Log(ui.ErrorWriter.String())
				t.Fatalf("expected %d to be %d", code, exp)
			}

			bundlePath := filepath.Join(testDir, basePath+debugCompressionExt)
			_, err = os.Open(bundlePath)
			if err != nil {
				t.Fatalf("failed to open archive: %s", err)
			}

			tgz := archiver.NewTarGz()
			err = tgz.Walk(bundlePath, func(f archiver.File) error {
				fh, ok := f.Header.(*tar.Header)
				if !ok {
					t.Fatalf("invalid file header: %#v", f.Header)
				}

				// Ignore base directory and index file
				if fh.Name == basePath+"/" || fh.Name == filepath.Join(basePath, "index.json") {
					return nil
				}

				for _, fileName := range tc.expectedFiles {
					if fh.Name == filepath.Join(basePath, fileName) {
						return nil
					}
				}

				// If we reach here, it means that this is an unexpected file
				return fmt.Errorf("unexpected file: %s", fh.Name)
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDebugCommand_Pprof(t *testing.T) {
	testDir, err := ioutil.TempDir("", "vault-debug")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	client, closer := testVaultServer(t)
	defer closer()

	ui, cmd := testDebugCommand(t)
	cmd.client = client
	cmd.skipTimingChecks = true

	basePath := "pprof"
	outputPath := filepath.Join(testDir, basePath)
	// pprof requires a minimum interval of 1s, we set it to 2 to ensure it
	// runs through and reduce flakiness on slower systems.
	args := []string{
		"-compress=false",
		"-duration=2s",
		"-interval=2s",
		fmt.Sprintf("-output=%s", outputPath),
		"-target=pprof",
	}

	code := cmd.Run(args)
	if exp := 0; code != exp {
		t.Log(ui.ErrorWriter.String())
		t.Fatalf("expected %d to be %d", code, exp)
	}

	profiles := []string{"heap.prof", "goroutine.prof"}
	pollingProfiles := []string{"profile.prof", "trace.out"}

	// These are captures on the first (0th) and last (1st) frame
	for _, v := range profiles {
		files, _ := filepath.Glob(fmt.Sprintf("%s/*/%s", outputPath, v))
		if len(files) != 2 {
			t.Errorf("2 output files should exist for %s: got: %v", v, files)
		}
	}

	// Since profile and trace are polling outputs, these only get captured
	// on the first (0th) frame.
	for _, v := range pollingProfiles {
		files, _ := filepath.Glob(fmt.Sprintf("%s/*/%s", outputPath, v))
		if len(files) != 1 {
			t.Errorf("1 output file should exist for %s: got: %v", v, files)
		}
	}

	t.Log(ui.OutputWriter.String())
	t.Log(ui.ErrorWriter.String())
}

func TestDebugCommand_IndexFile(t *testing.T) {
	t.Parallel()

	testDir, err := ioutil.TempDir("", "vault-debug")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	client, closer := testVaultServer(t)
	defer closer()

	ui, cmd := testDebugCommand(t)
	cmd.client = client
	cmd.skipTimingChecks = true

	basePath := "index-test"
	outputPath := filepath.Join(testDir, basePath)
	// pprof requires a minimum interval of 1s
	args := []string{
		"-compress=false",
		"-duration=1s",
		"-interval=1s",
		"-metrics-interval=1s",
		fmt.Sprintf("-output=%s", outputPath),
	}

	code := cmd.Run(args)
	if exp := 0; code != exp {
		t.Log(ui.ErrorWriter.String())
		t.Fatalf("expected %d to be %d", code, exp)
	}

	content, err := ioutil.ReadFile(filepath.Join(outputPath, "index.json"))
	if err != nil {
		t.Fatal(err)
	}

	index := &debugIndex{}
	if err := json.Unmarshal(content, index); err != nil {
		t.Fatal(err)
	}
	if len(index.Output) == 0 {
		t.Fatalf("expected valid index file: got: %v", index)
	}
}

func TestDebugCommand_TimingChecks(t *testing.T) {
	t.Parallel()

	testDir, err := ioutil.TempDir("", "vault-debug")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	cases := []struct {
		name            string
		duration        string
		interval        string
		metricsInterval string
	}{
		{
			"short-values-all",
			"10ms",
			"10ms",
			"10ms",
		},
		{
			"short-duration",
			"10ms",
			"",
			"",
		},
		{
			"short-interval",
			debugMinInterval.String(),
			"10ms",
			"",
		},
		{
			"short-metrics-interval",
			debugMinInterval.String(),
			"",
			"10ms",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			// If we are past the minimum duration + some grace, trigger shutdown
			// to prevent hanging
			grace := 10 * time.Second
			shutdownCh := make(chan struct{})
			go func() {
				time.AfterFunc(grace, func() {
					close(shutdownCh)
				})
			}()

			ui, cmd := testDebugCommand(t)
			cmd.client = client
			cmd.ShutdownCh = shutdownCh

			basePath := tc.name
			outputPath := filepath.Join(testDir, basePath)
			// pprof requires a minimum interval of 1s
			args := []string{
				"-target=server-status",
				fmt.Sprintf("-output=%s", outputPath),
			}
			if tc.duration != "" {
				args = append(args, fmt.Sprintf("-duration=%s", tc.duration))
			}
			if tc.interval != "" {
				args = append(args, fmt.Sprintf("-interval=%s", tc.interval))
			}
			if tc.metricsInterval != "" {
				args = append(args, fmt.Sprintf("-metrics-interval=%s", tc.metricsInterval))
			}

			code := cmd.Run(args)
			if exp := 0; code != exp {
				t.Log(ui.ErrorWriter.String())
				t.Fatalf("expected %d to be %d", code, exp)
			}

			if !strings.Contains(ui.OutputWriter.String(), "Duration: 5s") {
				t.Fatal("expected minimum duration value")
			}

			if tc.interval != "" {
				if !strings.Contains(ui.OutputWriter.String(), "  Interval: 5s") {
					t.Fatal("expected minimum interval value")
				}
			}

			if tc.metricsInterval != "" {
				if !strings.Contains(ui.OutputWriter.String(), "Metrics Interval: 5s") {
					t.Fatal("expected minimum metrics interval value")
				}
			}
		})
	}
}

func TestDebugCommand_NoConnection(t *testing.T) {
	t.Parallel()

	client, err := api.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	_, cmd := testDebugCommand(t)
	cmd.client = client
	cmd.skipTimingChecks = true

	args := []string{
		"-duration=1s",
		"-target=server-status",
	}

	code := cmd.Run(args)
	if exp := 1; code != exp {
		t.Fatalf("expected %d to be %d", code, exp)
	}
}

func TestDebugCommand_OutputExists(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		compress      bool
		outputFile    string
		expectedError string
	}{
		{
			"no-compress",
			false,
			"output-exists",
			"output directory already exists",
		},
		{
			"compress",
			true,
			"output-exist.tar.gz",
			"output file already exists",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testDir, err := ioutil.TempDir("", "vault-debug")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(testDir)

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testDebugCommand(t)
			cmd.client = client
			cmd.skipTimingChecks = true

			outputPath := filepath.Join(testDir, tc.outputFile)

			// Create a conflicting file/directory
			if tc.compress {
				_, err = os.Create(outputPath)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				err = os.Mkdir(outputPath, 0o755)
				if err != nil {
					t.Fatal(err)
				}
			}

			args := []string{
				fmt.Sprintf("-compress=%t", tc.compress),
				"-duration=1s",
				"-interval=1s",
				"-metrics-interval=1s",
				fmt.Sprintf("-output=%s", outputPath),
			}

			code := cmd.Run(args)
			if exp := 1; code != exp {
				t.Log(ui.OutputWriter.String())
				t.Log(ui.ErrorWriter.String())
				t.Errorf("expected %d to be %d", code, exp)
			}

			output := ui.ErrorWriter.String() + ui.OutputWriter.String()
			if !strings.Contains(output, tc.expectedError) {
				t.Fatalf("expected %s, got: %s", tc.expectedError, output)
			}
		})
	}
}

func TestDebugCommand_PartialPermissions(t *testing.T) {
	t.Parallel()

	testDir, err := ioutil.TempDir("", "vault-debug")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	client, closer := testVaultServer(t)
	defer closer()

	// Create a new token with default policy
	resp, err := client.Logical().Write("auth/token/create", map[string]interface{}{
		"policies": "default",
	})
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken(resp.Auth.ClientToken)

	ui, cmd := testDebugCommand(t)
	cmd.client = client
	cmd.skipTimingChecks = true

	basePath := "with-default-policy-token"
	args := []string{
		"-duration=1s",
		fmt.Sprintf("-output=%s/%s", testDir, basePath),
	}

	code := cmd.Run(args)
	if exp := 0; code != exp {
		t.Log(ui.ErrorWriter.String())
		t.Fatalf("expected %d to be %d", code, exp)
	}

	bundlePath := filepath.Join(testDir, basePath+debugCompressionExt)
	_, err = os.Open(bundlePath)
	if err != nil {
		t.Fatalf("failed to open archive: %s", err)
	}

	tgz := archiver.NewTarGz()
	err = tgz.Walk(bundlePath, func(f archiver.File) error {
		fh, ok := f.Header.(*tar.Header)
		if !ok {
			t.Fatalf("invalid file header: %#v", f.Header)
		}

		// Ignore base directory and index file
		if fh.Name == basePath+"/" {
			return nil
		}

		// Ignore directories, which still get created by pprof but should
		// otherwise be empty.
		if fh.FileInfo().IsDir() {
			return nil
		}

		switch {
		case fh.Name == filepath.Join(basePath, "index.json"):
		case fh.Name == filepath.Join(basePath, "replication_status.json"):
		case fh.Name == filepath.Join(basePath, "server_status.json"):
		case fh.Name == filepath.Join(basePath, "vault.log"):
		default:
			return fmt.Errorf("unexpected file: %s", fh.Name)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
