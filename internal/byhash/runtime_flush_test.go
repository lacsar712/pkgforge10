package byhash

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDumpInReleaseFlushesLargeBody guards against regressions of the
// bufio.Writer flush bug: a body larger than the default 4KiB buffer
// must be fully persisted to disk, not left stranded in the buffer when
// the file is closed.
func TestDumpInReleaseFlushesLargeBody(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dists", "stable", "InRelease")
	// ~14KB, well past bufio's default 4096-byte buffer.
	body := strings.Repeat("sha256=abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789\n", 200)
	if err := DumpInRelease(path, body); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != body {
		t.Fatalf("body not fully persisted: got %d bytes, want %d", len(b), len(body))
	}
}
