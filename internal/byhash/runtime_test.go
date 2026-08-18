package byhash

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExportPoolFileRejectsEscape(t *testing.T) {
	if _, err := ExportPoolFile(t.TempDir(), filepath.Join("..", "etc", "passwd")); err == nil {
		t.Fatal("expected path escape to be rejected")
	}
}

func TestParseReleaseJSONRejectsInvalid(t *testing.T) {
	if _, err := ParseReleaseJSON([]byte("SHA256: not-json")); err == nil {
		t.Fatal("expected JSON error")
	}
}

func TestNilDebName(t *testing.T) {
	var d *Deb
	if d.Name() != "" {
		t.Fatalf("got %q", d.Name())
	}
}

func TestHashBagPutGet(t *testing.T) {
	bag := NewHashBag()
	bag.Put("foo.deb", "abc")
	if bag.Get("foo.deb") != "abc" {
		t.Fatal("hash not stored")
	}
}

func TestGrowHashNoWriteThrough(t *testing.T) {
	dst := make([]byte, 2, 8)
	copy(dst, []byte("AB"))
	got := GrowHash(dst, 'C')
	got[0] = 'X'
	if dst[0] != 'A' {
		t.Fatal("GrowHash wrote through into the digest buffer")
	}
}

func TestWrapSuiteDeniedIs(t *testing.T) {
	err := WrapSuiteDenied("publish", "stable")
	if !errors.Is(err, ErrSuite) {
		t.Fatalf("lost ErrSuite: %v", err)
	}
}

func TestWaitPublishHonorsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()
	start := time.Now()
	err := WaitPublish(ctx, 600*time.Millisecond)
	if err == nil {
		t.Fatal("expected cancel error")
	}
	if time.Since(start) > 250*time.Millisecond {
		t.Fatalf("WaitPublish ignored cancel, elapsed=%s", time.Since(start))
	}
}

func TestCopySHAIndependent(t *testing.T) {
	src := []byte("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	got := CopySHA(src, 8)
	got[0] = 'f'
	if src[0] != '0' {
		t.Fatal("CopySHA aliased the sha256 hex")
	}
}

func TestAfterWriteRejectsEpochRollback(t *testing.T) {
	min := ""
	get := func() (string, error) { return min, nil }
	set := func(v string) error { min = v; return nil }
	body := Sample().Body + " epoch=5"
	if err := AfterWrite(get, set, body); err != nil {
		t.Fatal(err)
	}
	if err := AfterWrite(get, set, Sample().Body+" epoch=3"); err == nil {
		t.Fatal("expected epoch rollback to be rejected")
	}
}

func TestDumpInReleasePersists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "InRelease")
	body := "Origin: pkgforge\n"
	if err := DumpInRelease(path, body); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != body {
		t.Fatalf("got %q", b)
	}
}
