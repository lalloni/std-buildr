package ar

import (
	"archive/tar"
	"io"
	"os"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/ulikunitz/xz"
)

func TarXz(target string, files []string, namer func(name string) string) error {
	w, err := os.Create(target)
	if err != nil {
		return errors.Wrapf(err, "creating target package file")
	}
	defer w.Close()
	xzw, err := xz.NewWriter(w)
	if err != nil {
		return errors.Wrapf(err, "creating target xz stream writer")
	}
	defer xzw.Close()
	tw := tar.NewWriter(xzw)
	for _, file := range files {
		log.Infof("packaging processed file '%s'", file)
		fi, err := os.Stat(file)
		if err != nil {
			return errors.Wrapf(err, "reading info of processed file '%s'", file)
		}
		h, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return errors.Wrapf(err, "building tar header for processes file '%s'", file)
		}
		h.Name = namer(file)
		h.Uid = 0
		h.Gid = 0
		err = tw.WriteHeader(h)
		if err != nil {
			return errors.Wrapf(err, "writing tar header for processed file '%s'", file)
		}
		in, err := os.Open(file)
		if err != nil {
			return errors.Wrapf(err, "opening processed file '%s'", file)
		}
		_, err = io.Copy(tw, in)
		in.Close()
		if err != nil {
			return errors.Wrapf(err, "copying processed file '%s' data", file)
		}
	}
	err = tw.Close()
	if err != nil {
		return errors.Wrapf(err, "writing package file '%s'", target)
	}
	return nil
}
