package cli

import (
	"io"
	"os"
	"strings"
	"testing"

	"iasi-dev/internal/consts/RC"
	"iasi-dev/internal/structures"
)

func TestErrorMessageIsBoldRed(t *testing.T) {
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer read.Close()

	parms := structures.Parms{Verbose: int(visibilityNormal)}
	writeMessage(parms, write, visibilityNormal, levelError, true, "boom")
	if err := write.Close(); err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(read)
	if err != nil {
		t.Fatal(err)
	}
	output := string(data)

	if !strings.Contains(output, colorBold+colorRed+"boom"+colorReset) {
		t.Fatalf("error output is not bold red: %q", output)
	}
}

func TestHeaderWritesTimestampedBannerToLog(t *testing.T) {
	log, err := os.CreateTemp(t.TempDir(), "iasi-log-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()

	parms := structures.Parms{Verbose: 0, LogFile: log}
	Header(parms, "Releasing %s", "iasi-quarto")

	if err := log.Sync(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(log.Name())
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("header log has %d lines, want 3: %q", len(lines), string(data))
	}

	for i, line := range lines {
		if len(line) < 11 || line[2] != ':' || line[5] != ':' || line[8:11] != " - " {
			t.Fatalf("line %d has no timestamp prefix: %q", i+1, line)
		}
	}

	if !strings.HasSuffix(lines[0], logBanner) {
		t.Fatalf("first banner line missing: %q", lines[0])
	}
	if !strings.HasSuffix(lines[1], "Releasing iasi-quarto") {
		t.Fatalf("header line missing: %q", lines[1])
	}
	if !strings.HasSuffix(lines[2], logBanner) {
		t.Fatalf("last banner line missing: %q", lines[2])
	}
}

func TestHeaderMessageIsBoldCyan(t *testing.T) {
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer read.Close()

	message := "Releasing iasi-quarto"
	writeConsoleMessage(write, levelHeader, true, message)
	if err := write.Close(); err != nil {
		t.Fatal(err)
	}

	data, err := io.ReadAll(read)
	if err != nil {
		t.Fatal(err)
	}
	output := string(data)

	if !strings.Contains(output, colorBold+colorCyan+message+colorReset) {
		t.Fatalf("header output is not bold cyan: %q", output)
	}
}

func TestErrorMessageDoesNotAbort(t *testing.T) {
	rc := RC.OK
	parms := structures.Parms{RC: &rc, Verbose: int(visibilityNormal)}

	completed := false
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				t.Fatalf("ErrorMessage aborted: %v", recovered)
			}
		}()
		ErrorMessage(parms, "rollback warning")
		completed = true
	}()

	if !completed {
		t.Fatal("ErrorMessage did not return normally")
	}
	if rc != RC.OK {
		t.Fatalf("rc = 0x%02X, want 0x%02X", rc, RC.OK)
	}
}

func TestErrorStopsAfterReporting(t *testing.T) {
	rc := RC.OK
	parms := structures.Parms{RC: &rc, Verbose: int(visibilityNormal)}

	defer func() {
		recovered := recover()
		stop, ok := recovered.(RC.Stop)
		if !ok {
			t.Fatalf("Error recovered %T, want RC.Stop", recovered)
		}
		if stop.Code != RC.Error {
			t.Fatalf("stop code = 0x%02X, want 0x%02X", stop.Code, RC.Error)
		}
		if rc != RC.Error {
			t.Fatalf("rc = 0x%02X, want 0x%02X", rc, RC.Error)
		}
	}()

	Error(RC.Error, parms, "boom")
}

func TestAbortStopsWithoutAdditionalMessage(t *testing.T) {
	rc := RC.OK
	parms := structures.Parms{RC: &rc}

	defer func() {
		recovered := recover()
		stop, ok := recovered.(RC.Stop)
		if !ok {
			t.Fatalf("Abort recovered %T, want RC.Stop", recovered)
		}
		if stop.Code != RC.Error {
			t.Fatalf("stop code = 0x%02X, want 0x%02X", stop.Code, RC.Error)
		}
		if rc != RC.Error {
			t.Fatalf("rc = 0x%02X, want 0x%02X", rc, RC.Error)
		}
	}()

	Abort(RC.Error, parms)
}
