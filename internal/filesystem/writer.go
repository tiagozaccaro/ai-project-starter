// Package filesystem provides safe file writing with atomic
// write-then-rename semantics, dry-run collection, and
// destination safety checks.
package filesystem

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Writer writes rendered files to a destination directory.
type Writer struct {
	Destination string
	DryRun      bool
}

// NewWriter returns a Writer for the destination directory.
func NewWriter(destination string, dryRun bool) *Writer {
	return &Writer{Destination: destination, DryRun: dryRun}
}

// WriteResult reports the outcome of a write operation.
type WriteResult struct {
	Written   []string
	DryRun    bool
	FileCount int
}

// WriteFiles writes the given rendered files to the
// destination. In dry-run mode nothing is written; the
// planned file list is returned instead.
func (w *Writer) WriteFiles(files []File) (*WriteResult, error) {
	result := &WriteResult{DryRun: w.DryRun}
	if w.DryRun {
		for _, f := range files {
			result.Written = append(result.Written, f.TargetPath)
		}
		result.FileCount = len(files)
		return result, nil
	}

	if err := os.MkdirAll(w.Destination, 0o755); err != nil {
		return nil, fmt.Errorf("create destination: %w", err)
	}

	for _, f := range files {
		if err := w.writeFile(f); err != nil {
			return nil, err
		}
		result.Written = append(result.Written, f.TargetPath)
	}
	result.FileCount = len(result.Written)
	return result, nil
}

// writeFile atomically writes a single file via a temporary
// file and rename, creating parent directories as needed.
func (w *Writer) writeFile(f File) error {
	target := filepath.Join(w.Destination, f.TargetPath)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("create dir for %s: %w", f.TargetPath, err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(target), ".tmp-*")
	if err != nil {
		return fmt.Errorf("create temp for %s: %w", f.TargetPath, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := bytes.NewReader(f.Content).WriteTo(tmp); err != nil {
		tmp.Close()
		return fmt.Errorf("write %s: %w", f.TargetPath, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close %s: %w", f.TargetPath, err)
	}
	mode := os.FileMode(0o644)
	if f.Executable {
		mode = 0o755
	}
	if err := os.Chmod(tmpName, mode); err != nil {
		return fmt.Errorf("chmod %s: %w", f.TargetPath, err)
	}
	if err := os.Rename(tmpName, target); err != nil {
		return fmt.Errorf("rename %s: %w", f.TargetPath, err)
	}
	return nil
}

// File is a single output file.
type File struct {
	TargetPath string
	Content    []byte
	Executable bool
}

// DestinationEmpty reports whether the destination directory
// is empty or does not exist.
func DestinationEmpty(destination string) (bool, error) {
	info, err := os.Stat(destination)
	if os.IsNotExist(err) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	if !info.IsDir() {
		return false, fmt.Errorf("destination %s is not a directory", destination)
	}
	entries, err := os.ReadDir(destination)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

// ListDestinationFiles returns the relative paths of all
// files currently in the destination directory.
func ListDestinationFiles(destination string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(destination, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(destination, path)
		if err != nil {
			return err
		}
		out = append(out, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
