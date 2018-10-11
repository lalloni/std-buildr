package ar

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/ulikunitz/xz"
)

var UnknownFormat = errors.New("unknown format")

type format string

var (
	TarXzFormat format = "tar.xz"
	TarGzFormat format = "tar.gz"
	ZipFormat   format = "zip"
)

func (f *format) AddExt(name string) string {
	return name + "." + string(*f)
}

func (f *format) ChangeExt(name string) string {
	return f.AddExt(strings.TrimSuffix(name, filepath.Ext(name)))
}

func Format(s string) (format, error) {
	switch s {
	case string(TarGzFormat):
		return TarGzFormat, nil
	case string(TarXzFormat):
		return TarXzFormat, nil
	case string(ZipFormat):
		return ZipFormat, nil
	default:
		return "", UnknownFormat
	}
}

func FormatDefault(s string, def format) (format, error) {
	if s == "" {
		return def, nil
	}
	return Format(s)
}

func Package(targetFormat format, target string, files []string) error {
	var err error
	switch targetFormat {
	case TarXzFormat:
		err = TarXz(target, files, path.Base)
	case TarGzFormat:
		err = TarGz(target, files, path.Base)
	case ZipFormat:
		err = Zip(target, files, path.Base)
	default:
		return UnknownFormat
	}
	if err != nil {
		return errors.Wrapf(err, "packaging files in %s", targetFormat)
	}
	return nil
}

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

func TarGz(target string, files []string, namer func(name string) string) error {
	w, err := os.Create(target)
	if err != nil {
		return errors.Wrapf(err, "creating target package file")
	}
	defer w.Close()
	gzw := gzip.NewWriter(w)
	if err != nil {
		return errors.Wrapf(err, "creating target gz stream writer")
	}
	defer gzw.Close()
	tw := tar.NewWriter(gzw)
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

func Zip(target string, files []string, namer func(name string) string) error {
	w, err := os.Create(target)
	if err != nil {
		return errors.Wrapf(err, "creating target package file")
	}
	defer w.Close()
	zipw := zip.NewWriter(w)
	defer zipw.Close()
	for _, file := range files {
		log.Infof("packaging processed file '%s'", file)
		in, err := os.Open(file)
		if err != nil {
			return errors.Wrapf(err, "opening processed file '%s'", file)
		}
		defer in.Close()
		info, err := in.Stat()
		if err != nil {
			return err
		}
		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return errors.Wrapf(err, "getting info for processed file '%s'", file)
		}
		h.Name = namer(file)
		h.Method = zip.Deflate
		writer, err := zipw.CreateHeader(h)
		if err != nil {
			return errors.Wrapf(err, "writing zip header for processed file '%s'", file)
		}
		if _, err = io.Copy(writer, in); err != nil {
			return errors.Wrapf(err, "copying processed file '%s' data", file)
		}
		in.Close()
	}
	err = zipw.Close()
	if err != nil {
		return errors.Wrapf(err, "writing package file '%s'", target)
	}
	return nil

}
