package sh

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/apex/log"
)

// Rm removes the given file or directory even if non-empty. It will not return
// an error if the target doesn't exist, only if the target cannot be removed.
func Rm(path string) error {
	err := os.RemoveAll(path)
	if err == nil || os.IsNotExist(err) {
		return nil
	}
	return fmt.Errorf(`failed to remove %s: %v`, path, err)
}

func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, "checking path existence")
	}
	return true, nil
}

func FirstExist(paths ...string) (string, error) {
	for _, path := range paths {
		b, err := Exist(path)
		if err != nil {
			return "", errors.Wrap(err, "searching first existent")
		}
		if b {
			return path, nil
		}
	}
	return "", os.ErrNotExist
}

func All(os.FileInfo) bool {
	return true
}

func Files(fi os.FileInfo) bool {
	return !fi.IsDir()
}

func Collect(path string, accept func(os.FileInfo) bool) ([]string, error) {
	r := []string{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if accept(info) {
			r = append(r, path)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "collecting contents of %s", path)
	}
	return r, nil
}

func CollectFiles(path string) ([]string, error) {
	return Collect(path, Files)
}

func newLogWriter(entry *log.Entry, prefix string) *logWriter {
	return &logWriter{
		entry:  entry,
		prefix: prefix,
	}
}

type logWriter struct {
	entry  *log.Entry
	prefix string
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	w.entry.Infof("%s: %s", w.prefix, string(p))
	return len(p), nil
}
