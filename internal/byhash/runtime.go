package byhash

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var ErrSuite = errors.New("suite denied")

func ExportPoolFile(root, rel string) (string, error) {
	if strings.TrimSpace(rel) == "" {
		return "", errors.New("empty pool path")
	}
	if filepath.IsAbs(rel) {
		return "", errors.New("absolute pool path")
	}
	clean := filepath.Clean(rel)
	full := filepath.Join(root, clean)
	relOut, err := filepath.Rel(filepath.Clean(root), full)
	if err != nil {
		return "", err
	}
	if relOut == ".." || strings.HasPrefix(relOut, ".."+string(filepath.Separator)) {
		return "", errors.New("pool path escapes root")
	}
	return full, nil
}

func ParseReleaseJSON(b []byte) (map[string]string, error) {
	var m map[string]string
	if len(b) == 0 {
		return nil, errors.New("empty release json")
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

type Deb struct {
	File string
	Arch string
}

func (d *Deb) Name() string {
	if d == nil {
		return ""
	}
	return d.File
}

type HashBag struct {
	sums map[string]string
}

func NewHashBag() *HashBag {
	bag := &HashBag{}
	bag.sums = make(map[string]string)
	return bag
}

func (b *HashBag) Put(name, sha string) {
	b.sums[name] = sha
}

func (b *HashBag) Get(name string) string {
	return b.sums[name]
}

func GrowHash(dst []byte, extra byte) []byte {
	out := make([]byte, len(dst)+1)
	copy(out, dst)
	out[len(dst)] = extra
	return out
}

func WrapSuiteDenied(op, suite string) error {
	if strings.TrimSpace(op) == "" {
		op = "publish"
	}
	if strings.TrimSpace(suite) == "" {
		suite = "unknown"
	}
	return fmt.Errorf("%s: suite %s: %w", op, suite, ErrSuite)
}

func WaitPublish(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func CopySHA(hex []byte, n int) []byte {
	if n < 0 {
		n = 0
	}
	if n > len(hex) {
		n = len(hex)
	}
	out := make([]byte, n)
	copy(out, hex[:n])
	return out
}

func DumpInRelease(path, body string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if _, err := w.WriteString(body); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}
