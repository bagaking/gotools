package fop

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFileTruncatesExistingDestination(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dest := filepath.Join(dir, "dest.txt")

	if err := os.WriteFile(src, []byte("new"), 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", src, err)
	}
	if err := os.WriteFile(dest, []byte("old trailing data"), 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", dest, err)
	}

	if err := CopyFile(src, dest, false); err != nil {
		t.Fatalf("CopyFile(%q, %q, false) error = %v", src, dest, err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error = %v", dest, err)
	}
	if string(got) != "new" {
		t.Errorf("CopyFile(%q, %q, false) wrote %q, want %q", src, dest, string(got), "new")
	}
}

func TestCopyFileSamePathNoOp(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	want := []byte("keep source content")

	if err := os.WriteFile(src, want, 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", src, err)
	}

	if err := CopyFile(src, src, false); err != nil {
		t.Fatalf("CopyFile(%q, %q, false) error = %v", src, src, err)
	}

	got, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error = %v", src, err)
	}
	if string(got) != string(want) {
		t.Errorf("CopyFile(%q, %q, false) wrote %q, want %q", src, src, string(got), string(want))
	}
}

func TestCopyFileHardLinkNoOp(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dest := filepath.Join(dir, "dest.txt")
	want := []byte("keep linked content")

	if err := os.WriteFile(src, want, 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", src, err)
	}
	if err := os.Link(src, dest); err != nil {
		t.Fatalf("os.Link(%q, %q) error = %v", src, dest, err)
	}

	if err := CopyFile(src, dest, false); err != nil {
		t.Fatalf("CopyFile(%q, %q, false) error = %v", src, dest, err)
	}

	got, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error = %v", src, err)
	}
	if string(got) != string(want) {
		t.Errorf("CopyFile(%q, %q, false) changed source to %q, want %q", src, dest, string(got), string(want))
	}
}
