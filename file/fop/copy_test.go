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
