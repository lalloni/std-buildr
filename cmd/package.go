// Copyright © 2018 Pablo Lalloni <plalloni@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/Masterminds/semver"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ulikunitz/xz"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/sh"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package the project",
	RunE:  runPackage,
}

var (
	tagNameRegexp     = regexp.MustCompile(`^v(.*)$`)
	includeRegexp     = regexp.MustCompile(`^@@(.*)$`)
	evolutionalRegexp = regexp.MustCompile(`^(.+-)?([0-9]{6,})-(dml|dcl|ddl)(-.+)?\.sql$`)
)

func init() {
	rootCmd.AddCommand(packageCmd)
}

type packageConfig struct {
	Version    *semver.Version
	Untracked  bool
	Changed    bool
	Uncommited bool
}

func (p *packageConfig) Dirty() bool {
	return p.Untracked || p.Uncommited || p.Changed
}

func runPackage(cmd *cobra.Command, args []string) error {
	gitversion, err := sh.Output("git", "--version")
	if err != nil {
		return errors.Wrap(err, "determining git version")
	}
	log.Info(gitversion)

	packageCfg := &packageConfig{}

	v1, err := sh.Output("git", "describe", "--abbrev=40", "HEAD")
	if err != nil {
		return errors.Wrap(err, "getting version from git")
	}
	v2 := tagNameRegexp.FindStringSubmatch(v1)
	if v2 == nil {
		return errors.Errorf("tag name must be prefixed with a 'v' character (found '%s')", packageCfg.Version)
	}
	packageCfg.Version, err = semver.NewVersion(v2[1])
	if err != nil {
		return errors.Wrapf(err, "tag name must be a valid semver 2 string prefixed with a 'v' character (found '%s')", packageCfg.Version)
	}
	log.Infof("version is '%s'", packageCfg.Version)

	s, err := sh.Output("git", "ls-files", "--exclude-standard", "--others")
	if err != nil {
		return errors.Wrap(err, "listing untracked files from git")
	}
	packageCfg.Untracked = len(s) > 0
	log.Infof("untracked files present: %v", packageCfg.Untracked)

	err = sh.Run("git", "diff-files", "--quiet")
	if err != nil {
		packageCfg.Changed = true
	}
	log.Infof("changed tracked files present: %v", packageCfg.Changed)

	err = sh.Run("git", "diff-index", "--cached", "--quiet", "HEAD")
	if err != nil {
		packageCfg.Uncommited = true
	}
	log.Infof("uncommited staged files present: %v", packageCfg.Uncommited)

	cfg := viper.Get("buildr.config").(*config.Config)
	log.Infof("config is %+v", cfg)

	switch cfg.Type {
	case config.TypeOracleSQLEvolutional:
		err = packageOracleSQLEvolutional(cfg, packageCfg)
	default:
		err = errors.Errorf("type not implemented: '%s'", cfg.Type)
	}
	return err
}

func packageOracleSQLEvolutional(c *config.Config, p *packageConfig) error {
	const targetSource = "target/source"
	err := os.MkdirAll(targetSource, 0775)
	if err != nil {
		return errors.Wrapf(err, "creating directory")
	}
	base, err := sh.FirstExist("src/sql/inc", "src/sql/incremental")
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Errorf("missing incremental sources (at 'src/sql/inc[remental]')")
		}
		return errors.Wrapf(err, "checking incremental source presence")
	}

	// preprocess
	sources, err := sh.CollectFiles(base)
	if err != nil {
		return errors.Wrapf(err, "collecting source files from '%s'", base)
	}

	for _, source := range sources {

		if !evolutionalRegexp.MatchString(path.Base(source)) {
			return errors.Errorf("source file name '%s' does not match standard naming scheme (%s)", source, evolutionalRegexp.String())
		}

		ss := evolutionalRegexp.FindStringSubmatch(path.Base(source))
		if len(ss[1]) != 0 && ss[1] != c.ApplicationID+"-" {
			return errors.Errorf("source file '%s' name prefix '%s' must equal application id '%s' if used", source, ss[1][:len(ss[1])-1], c.ApplicationID)
		}

		targetName := path.Base(source)
		if len(ss[1]) == 0 {
			targetName = c.ApplicationID + "-" + targetName
		}

		log.Infof("processing source file '%s'", source)
		target := targetSource + "/" + targetName
		in, err := os.Open(source)
		if err != nil {
			return errors.Wrapf(err, "opening '%s'", source)
		}
		defer in.Close()
		log.Infof("into target file '%s'", target)
		out, err := os.Create(target)
		if err != nil {
			in.Close()
			return errors.Wrapf(err, "creating '%s'", target)
		}
		defer out.Close()
		s := bufio.NewScanner(in)
		for s.Scan() {
			l := s.Text()
			ms := includeRegexp.FindStringSubmatch(l)
			if ms == nil {
				_, err := fmt.Fprintln(out, l)
				if err != nil {
					return errors.Wrap(err, "copying input to output")
				}
			} else {
				i := filepath.Clean(filepath.Join(filepath.Dir(source), ms[1]))
				log.Infof("including '%s'", i)
				inc, err := os.Open(i)
				if err != nil {
					return errors.Wrapf(err, "opening '%s' include '%s'", source, ms[1])
				}
				fmt.Fprintln(out, "-- begin include "+ms[1])
				_, err = io.Copy(out, inc)
				inc.Close()
				if err != nil {
					return errors.Wrapf(err, "copying include contents from '%s' into '%s'", ms[1], source, target)
				}
				fmt.Fprintln(out, "-- end include "+ms[1])
			}
		}
		in.Close()
		out.Close()
	}

	// package
	sources, err = sh.CollectFiles(targetSource)
	if err != nil {
		return errors.Wrapf(err, "collecting preprocessed files from '%s'", targetSource)
	}
	targetPackage := fmt.Sprintf("target/%s-%s-%s.tar.xz", c.SystemID, c.ApplicationID, p.Version)
	log.Infof("writing to '%s'", targetPackage)
	w, err := os.Create(targetPackage)
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
	for _, source := range sources {
		n := path.Base(source)
		log.Infof("packaging processed file '%s'", source)
		fi, err := os.Stat(source)
		if err != nil {
			return errors.Wrapf(err, "reading info of processed file '%s'", source)
		}
		h, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return errors.Wrapf(err, "building tar header for processes file '%s'", source)
		}
		h.Name = n
		err = tw.WriteHeader(h)
		if err != nil {
			return errors.Wrapf(err, "writing tar header for processed file '%s'", source)
		}
		in, err := os.Open(source)
		if err != nil {
			return errors.Wrapf(err, "opening processed file '%s'", source)
		}
		_, err = io.Copy(tw, in)
		in.Close()
		if err != nil {
			return errors.Wrapf(err, "copying processed file '%s' data", source)
		}
	}
	err = tw.Close()
	if err != nil {
		return errors.Wrapf(err, "writing package file '%s'", targetPackage)
	}
	log.Info("done")
	return nil
}
