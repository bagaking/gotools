package fscan

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewAndScanFilterIgnoreDirectoryByName(t *testing.T) {
	tests := []struct {
		name      string
		filter    FileFilter
		wantPaths []string
	}{
		{
			name:   "hidden directory",
			filter: FilterIgnoreHiddenFile(),
			wantPaths: []string{
				".",
				"skip",
				"skip/ignored.txt",
				"visible.txt",
			},
		},
		{
			name:   "named directory",
			filter: FilterIgnoreNameTableOfFile("skip"),
			wantPaths: []string{
				".",
				".hidden",
				".hidden/hidden.txt",
				"visible.txt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := makeScanFixture(t)
			got, err := NewAndScan(root, true, tt.filter)
			if err != nil {
				t.Fatalf("NewAndScan(%q, true, %s) error = %v, want nil", root, tt.name, err)
			}

			if !reflect.DeepEqual(got.GetPaths(), tt.wantPaths) {
				t.Errorf("NewAndScan(%q, true, %s).GetPaths() = %v, want %v", root, tt.name, got.GetPaths(), tt.wantPaths)
			}
		})
	}
}

func makeScanFixture(t *testing.T) string {
	t.Helper()

	root, err := ioutil.TempDir("", "gotools-fscan-")
	if err != nil {
		t.Fatalf("TempDir() error = %v, want nil", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(root); err != nil {
			t.Errorf("RemoveAll(%q) error = %v, want nil", root, err)
		}
	})

	writeFile(t, filepath.Join(root, "visible.txt"), "visible")
	writeFile(t, filepath.Join(root, ".hidden", "hidden.txt"), "hidden")
	writeFile(t, filepath.Join(root, "skip", "ignored.txt"), "ignored")

	return root
}

func writeFile(t *testing.T, path string, data string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v, want nil", filepath.Dir(path), err)
	}
	if err := ioutil.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v, want nil", path, err)
	}
}
