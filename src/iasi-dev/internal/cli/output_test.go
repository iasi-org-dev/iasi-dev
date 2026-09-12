package cli

import (
	"io"
	"os"
	"strings"
	"testing"

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
