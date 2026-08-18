package byhash

import (
	"fmt"
	"strconv"
	"strings"
)

type Rec struct {
	Title, Body string
	Tags        []string
}

func Sample() Rec {
	return Rec{
		Title: "foo_1.0.0_amd64.deb",
		Body:  "sha256=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		Tags:  []string{"stable", "main"},
	}
}

func Seed() []Rec {
	return []Rec{
		Sample(),
		{
			Title: "InRelease",
			Body:  "sha256=abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
			Tags:  []string{"stable"},
		},
	}
}

func AfterWrite(getMin func() (string, error), setMin func(string) error, body string) error {
	epoch := parseEpoch(body)
	if epoch < 0 {
		return fmt.Errorf("epoch missing in %q", body)
	}
	cur, err := getMin()
	if err == nil && strings.TrimSpace(cur) != "" {
		last, conv := strconv.Atoi(strings.TrimSpace(cur))
		if conv == nil && epoch < last {
			return fmt.Errorf("epoch %d < last published %d", epoch, last)
		}
	}
	return setMin(strconv.Itoa(epoch))
}

func parseEpoch(body string) int {
	for _, part := range strings.Fields(body) {
		k, v, ok := strings.Cut(part, "=")
		if !ok || k != "epoch" {
			continue
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			return -1
		}
		return n
	}
	return 0
}

func Steps() []string { return []string{"name-check", "index-byhash", "export-dist"} }

func Enforce(title, body string, tags []string) error {
	if !validName(title) {
		return fmt.Errorf("title must be name_version_arch.deb or InRelease")
	}
	if !validSHA(body) {
		return fmt.Errorf("body must contain sha256= and 64 hex chars")
	}
	if len(tags) == 0 {
		return fmt.Errorf("suite/component tag required")
	}
	return nil
}

func validName(s string) bool {
	s = strings.TrimSpace(s)
	if s == "InRelease" {
		return true
	}
	return strings.HasSuffix(s, ".deb") && strings.Count(s, "_") >= 2
}

func validSHA(body string) bool {
	for _, part := range strings.Fields(body) {
		k, v, ok := strings.Cut(part, "=")
		if !ok || k != "sha256" {
			continue
		}
		if len(v) != 64 {
			return false
		}
		for _, r := range v {
			ok := (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
			if !ok {
				return false
			}
		}
		return true
	}
	return false
}
