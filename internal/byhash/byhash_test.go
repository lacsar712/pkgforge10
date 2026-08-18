package byhash

import "testing"

func TestAcceptDeb(t *testing.T) {
	if err := Enforce(Sample().Title, Sample().Body, Sample().Tags); err != nil {
		t.Fatal(err)
	}
}

func TestRejectName(t *testing.T) {
	if err := Enforce("foo.tgz", Sample().Body, []string{"stable"}); err == nil {
		t.Fatal("expected reject")
	}
}
