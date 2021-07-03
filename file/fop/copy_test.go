package fop

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bagaking/gotools/file/fpth"
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

func TestCopyFileWithLinkRemainEnsureDirCreatesParentForSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires extra privileges on windows")
	}

	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	src := filepath.Join(dir, "src-link")
	dest := filepath.Join(dir, "missing", "dest-link")
	if err := os.WriteFile(target, []byte("target"), 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", target, err)
	}
	if err := os.Symlink(target, src); err != nil {
		t.Fatalf("os.Symlink(%q, %q) error = %v", target, src, err)
	}

	if err := CopyFileWithLinkRemain(src, dest, true); err != nil {
		t.Fatalf("CopyFileWithLinkRemain(%q, %q, true) error = %v", src, dest, err)
	}

	got, err := os.Readlink(dest)
	if err != nil {
		t.Fatalf("os.Readlink(%q) error = %v", dest, err)
	}
	if got != target {
		t.Errorf("CopyFileWithLinkRemain(%q, %q, true) link target = %q, want %q", src, dest, got, target)
	}
}

func TestCopyDirReportsWalkErrWithoutPanicWhenContinuing(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dest := filepath.Join(dir, "dest")
	if err := os.Mkdir(src, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", src, err)
	}
	if err := os.Mkdir(dest, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", dest, err)
	}

	walkErr := errors.New("walk denied")
	oldWalk := fpth.Walk
	fpth.Walk = func(root string, fn filepath.WalkFunc) error {
		return fn(filepath.Join(root, "blocked"), nil, walkErr)
	}
	t.Cleanup(func() {
		fpth.Walk = oldWalk
	})

	err := CopyDir(src, dest, false, false)
	if err == nil {
		t.Fatalf("CopyDir(%q, %q, false, false) error = nil, want error", src, dest)
	}
	if !errors.Is(err, walkErr) {
		t.Errorf("CopyDir(%q, %q, false, false) error = %v, want wrapping %v", src, dest, err, walkErr)
	}
}

func TestCopyDirReportsNilFileInfoWithoutPanicWhenContinuing(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dest := filepath.Join(dir, "dest")
	if err := os.Mkdir(src, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", src, err)
	}
	if err := os.Mkdir(dest, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", dest, err)
	}

	oldWalk := fpth.Walk
	fpth.Walk = func(root string, fn filepath.WalkFunc) error {
		return fn(filepath.Join(root, "missing-info"), nil, nil)
	}
	t.Cleanup(func() {
		fpth.Walk = oldWalk
	})

	err := CopyDir(src, dest, false, false)
	if err == nil {
		t.Fatalf("CopyDir(%q, %q, false, false) error = nil, want error", src, dest)
	}
}

func TestCopyDirMapsPathsByRelativePath(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dest := filepath.Join(dir, "dest")
	nested := filepath.Join(src, "nested", "src", "file.txt")
	if err := os.MkdirAll(filepath.Dir(nested), 0755); err != nil {
		t.Fatalf("os.MkdirAll(%q) error = %v", filepath.Dir(nested), err)
	}
	if err := os.Mkdir(dest, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", dest, err)
	}
	if err := os.WriteFile(nested, []byte("content"), 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", nested, err)
	}

	if err := CopyDir(src, dest, false, true); err != nil {
		t.Fatalf("CopyDir(%q, %q, false, true) error = %v", src, dest, err)
	}

	wantPath := filepath.Join(dest, "nested", "src", "file.txt")
	got, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error = %v", wantPath, err)
	}
	if string(got) != "content" {
		t.Errorf("CopyDir(%q, %q, false, true) copied %q, want %q", src, dest, string(got), "content")
	}
}

func TestCopyDirContinuesAndReturnsCopyErrors(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod permissions are not enforced the same way on windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("permission bits do not make files unreadable for root")
	}

	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dest := filepath.Join(dir, "dest")
	unreadable := filepath.Join(src, "a-unreadable.txt")
	readable := filepath.Join(src, "z-readable.txt")
	if err := os.Mkdir(src, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", src, err)
	}
	if err := os.Mkdir(dest, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", dest, err)
	}
	if err := os.WriteFile(unreadable, []byte("blocked"), 0600); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", unreadable, err)
	}
	if err := os.WriteFile(readable, []byte("copied"), 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", readable, err)
	}
	if err := os.Chmod(unreadable, 0000); err != nil {
		t.Fatalf("os.Chmod(%q, 0000) error = %v", unreadable, err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(unreadable, 0600)
	})

	err := CopyDir(src, dest, false, false)
	if err == nil {
		t.Fatalf("CopyDir(%q, %q, false, false) error = nil, want error", src, dest)
	}
	got, readErr := os.ReadFile(filepath.Join(dest, "z-readable.txt"))
	if readErr != nil {
		t.Fatalf("os.ReadFile(%q) error = %v", filepath.Join(dest, "z-readable.txt"), readErr)
	}
	if string(got) != "copied" {
		t.Errorf("CopyDir(%q, %q, false, false) copied readable file = %q, want %q", src, dest, string(got), "copied")
	}
}

func TestCopyDirReturnsDestinationWriteErrorWhenContinuing(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod permissions are not enforced the same way on windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("permission bits do not make directories unwritable for root")
	}

	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dest := filepath.Join(dir, "dest")
	if err := os.Mkdir(src, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", src, err)
	}
	if err := os.Mkdir(dest, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", dest, err)
	}
	if err := os.WriteFile(filepath.Join(src, "file.txt"), []byte("content"), 0644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", filepath.Join(src, "file.txt"), err)
	}
	if err := os.Chmod(dest, 0555); err != nil {
		t.Fatalf("os.Chmod(%q, 0555) error = %v", dest, err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(dest, 0755)
	})

	err := CopyDir(src, dest, false, false)
	if err == nil {
		t.Fatalf("CopyDir(%q, %q, false, false) error = nil, want error", src, dest)
	}
}

func TestCopyDirAggregatesMultipleErrorsWhenContinuing(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dest := filepath.Join(dir, "dest")
	if err := os.Mkdir(src, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", src, err)
	}
	if err := os.Mkdir(dest, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", dest, err)
	}

	firstErr := errors.New("first")
	secondErr := errors.New("second")
	oldWalk := fpth.Walk
	fpth.Walk = func(root string, fn filepath.WalkFunc) error {
		if err := fn(filepath.Join(root, "first"), nil, firstErr); err != nil {
			return err
		}
		return fn(filepath.Join(root, "second"), nil, secondErr)
	}
	t.Cleanup(func() {
		fpth.Walk = oldWalk
	})

	err := CopyDir(src, dest, false, false)
	if err == nil {
		t.Fatalf("CopyDir(%q, %q, false, false) error = nil, want error", src, dest)
	}
	if !errors.Is(err, firstErr) {
		t.Errorf("CopyDir(%q, %q, false, false) error = %v, want wrapping %v", src, dest, err, firstErr)
	}
	if !errors.Is(err, secondErr) {
		t.Errorf("CopyDir(%q, %q, false, false) error = %v, want wrapping %v", src, dest, err, secondErr)
	}
}

func TestCopyDirReturnsWalkFunctionErrorWhenContinuing(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dest := filepath.Join(dir, "dest")
	if err := os.Mkdir(src, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", src, err)
	}
	if err := os.Mkdir(dest, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", dest, err)
	}

	walkErr := errors.New("walk failed")
	oldWalk := fpth.Walk
	fpth.Walk = func(root string, fn filepath.WalkFunc) error {
		return walkErr
	}
	t.Cleanup(func() {
		fpth.Walk = oldWalk
	})

	err := CopyDir(src, dest, false, false)
	if !errors.Is(err, walkErr) {
		t.Errorf("CopyDir(%q, %q, false, false) error = %v, want wrapping %v", src, dest, err, walkErr)
	}
}

func TestCopyDirStopsOnWalkErrorWhenRequested(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dest := filepath.Join(dir, "dest")
	if err := os.Mkdir(src, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", src, err)
	}
	if err := os.Mkdir(dest, 0755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v", dest, err)
	}

	walkErr := errors.New("walk denied")
	oldWalk := fpth.Walk
	fpth.Walk = func(root string, fn filepath.WalkFunc) error {
		if err := fn(filepath.Join(root, "blocked"), nil, walkErr); err != nil {
			return err
		}
		t.Fatal("CopyDir continued walking after errorStop=true")
		return nil
	}
	t.Cleanup(func() {
		fpth.Walk = oldWalk
	})

	err := CopyDir(src, dest, false, true)
	if !errors.Is(err, walkErr) {
		t.Errorf("CopyDir(%q, %q, false, true) error = %v, want wrapping %v", src, dest, err, walkErr)
	}
}
