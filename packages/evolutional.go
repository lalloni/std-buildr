package packages

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/ar"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/sh"
)

var (
	tagNameRegexp     = regexp.MustCompile(`^v(.*)$`)
	includeRegexp     = regexp.MustCompile(`^@@(.*)$`)
	evolutionalRegexp = regexp.MustCompile(`^(.+-)?([0-9]{6,})-(dml|dcl|ddl)(-.+)?\.sql$`)
)

func PackageOracleSQLEvolutional(ctx *context.Context, c *config.Config) error {

	err := VerifyStandardVersion(ctx)
	if err != nil {
		return errors.Wrapf(err, "verifying version")
	}

	const targetSource = "target/source"
	err = os.MkdirAll(targetSource, 0775)
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
	sourcesTargetMap := make(map[string]string)
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
		sourcesTargetMap[source] = target
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

	// package all

	err = PackageAllSQL(targetSource, ctx, c)
	if err != nil {
		return errors.Wrapf(err, "packaging all")
	}

	// package incrementals
	if len(c.From) > 0 {
		log.Infof("generating incremental packages")
		for _, from := range c.From {
			log.Infof("from %s: listing v%s tag files", from, from)
			include := make(map[string]string)
			for k, v := range sourcesTargetMap {
				include[k] = v
			}
			s, err := sh.Output("git", "ls-tree", "-r", "--name-only", "v"+from)
			if err != nil {
				return errors.Wrapf(err, "listing tag 'v%s' content: %s", from, s)
			}
			ss := bufio.NewScanner(bytes.NewBufferString(s))
			for ss.Scan() {
				f := ss.Text()
				delete(include, f)
				log.Infof("excluding %s", f)
			}
			targetSources := make([]string, 0)
			for _, v := range include {
				targetSources = append(targetSources, v)
			}
			targetPackage := fmt.Sprintf("target/%s-%s-from-%s.tar.xz", c.ApplicationID, ctx.Build.String(), from)
			log.Infof("writing to '%s'", targetPackage)
			err = ar.TarXz(targetPackage, targetSources, path.Base)
			if err != nil {
				return errors.Wrapf(err, "packaging source files")
			}
			ctx.AddArtifact(targetPackage)
		}
	}

	log.Info("done")
	return nil
}
