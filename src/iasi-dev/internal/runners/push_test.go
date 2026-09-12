package runners

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPushRepositoryUsesConfiguredOriginAndCurrentBranch(t *testing.T) {
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer read.Close()

	original := os.Stdout
	os.Stdout = write
	ok := pushRepository(filepath.Join("work", "repo-a"))
	if err := write.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = original

	if !ok {
		t.Fatal("pushRepository() = false, want true")
	}

	data, err := io.ReadAll(read)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(data)); got != "repo-a: git push -u origin HEAD" {
		t.Fatalf("pushRepository output = %q", got)
	}
}
