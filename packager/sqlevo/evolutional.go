package sqlevo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/ar"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/templates"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/msg"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/sh"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/version"
)

var (
	tagNameRegexp     = regexp.MustCompile(`^v(.*)$`)
	includeRegexp     = regexp.MustCompile(`^@@(.*)$`)
	evolutionalRegexp = regexp.MustCompile(`^(.+[-_])?([0-9]{6,})[-_](dml|dcl|ddl)([-_].+)?\.sql$`)
)

func Package(cfg *config.Config, ctx *context.Context) error {

	v := tagNameRegexp.FindStringSubmatch(ctx.Build.Version)
	if v == nil {
		return errors.Errorf("tag name must be prefixed with a 'v' character (found '%s')", ctx.Build.Version)
	}
	ctx.Build.Version = v[1]

	ev, err := version.ParseSemanticVersion(ctx.Build.Version)
	if err != nil {
		return errors.Wrap(err, "checking valid semantic version")
	}

	version := fmt.Sprintf("%d.%d.%d", ev.Major(), ev.Minor(), ev.Patch())

	allowDirty := viper.GetBool("buildr.allow-dirty")
	if ctx.Build.Dirty() && !allowDirty {
		untrackedAndChangedFiles, err := git.ListUntrackedFilesAndChangedFiles()
		if err != nil {
			return errors.Wrap(err, "checking untracked and changed files")
		}
		uu := strings.Join(untrackedAndChangedFiles, " ")
		if uu == "" {
			return errors.Errorf(msg.PACKAGE_COMMITED_ERROR, version, version, version)
		}
		return errors.Errorf(msg.PACKAGE_UNTRACKER_ERROR, uu, version, version, version)
	}

	allowUntagged := viper.GetBool("buildr.allow-untagged")
	if ctx.Build.Untagged() && !allowUntagged {
		return errors.Errorf(msg.PACKAGE_UNTAGGED_ERROR, version, version, version)
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

		if filepath.Base(source) == templates.README {
			continue
		}

		if !evolutionalRegexp.MatchString(filepath.Base(source)) {
			return errors.Errorf("source file name '%s' does not match standard naming scheme (%s)", source, evolutionalRegexp.String())
		}

		ss := evolutionalRegexp.FindStringSubmatch(filepath.Base(source))
		if len(ss[1]) != 0 && (ss[1] != cfg.ApplicationID+"-" || ss[1] != cfg.ApplicationID+"_") {
			return errors.Errorf("source file '%s' name prefix '%s' must equal application id '%s' if used", source, ss[1][:len(ss[1])-1], cfg.ApplicationID)
		}

		targetName := filepath.Base(source)
		if len(ss[1]) == 0 {
			targetName = cfg.ApplicationID + "-" + targetName
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
					return errors.Wrapf(err, "copying include contents from '%s' into '%s'", source, target)
				}
				fmt.Fprintln(out, "-- end include "+ms[1])
			}
		}
		in.Close()
		out.Close()
	}

	// package complete
	log.Infof("generating complete package")
	targetSources, err := sh.CollectFiles(targetSource)
	if err != nil {
		return errors.Wrapf(err, "collecting preprocessed files from '%s'", targetSource)
	}
	format, err := ar.FormatDefault(cfg.Package.Format, ar.ZipFormat)
	if err != nil {
		return errors.Wrapf(err, "invalid package format '%s' in configuration", cfg.Package.Format)
	}
	packageName := format.AddExt(fmt.Sprintf("%s-%s", cfg.ApplicationID, ctx.Build.String()))
	targetPackage := fmt.Sprintf("target/%s", packageName)
	log.Infof("writing to '%s'", targetPackage)
	err = ar.Package(format, targetPackage, targetSources)
	if err != nil {
		return errors.Wrapf(err, "packaging source files")
	}
	ctx.AddArtifact(packageName, targetPackage)

	// package incrementals
	if len(cfg.From) > 0 {

		log.Infof("generating incremental packages")
		for _, from := range cfg.From {

			log.Infof("from %s: listing v%s tag files", from, from)

			fromV := "v" + from
			if fromV == v[0] {
				fromV, err = GetPreviousTag(v[0])
				if err != nil {
					return errors.Wrapf(err, "getting previous tag")
				}
			}
			include := make([]string, 0)

			s, err := sh.Output("git", "diff", "--name-only", "--diff-filter=A", fromV, v[0])
			if err != nil {
				return errors.Wrapf(err, "listing tag 'v%s' content: %s", from, s)
			}
			ss := bufio.NewScanner(bytes.NewBufferString(s))
			for ss.Scan() {
				f := ss.Text()
				fileInBaseRegexp := regexp.MustCompile(`^src\/sql\/(?:inc|incremental)\/(.+\.sql)`)

				if !fileInBaseRegexp.MatchString(f) {

					log.Infof("excluding %s", f)
					continue
				}
				sss := evolutionalRegexp.FindStringSubmatch(filepath.Base(f))
				ff := filepath.Base(f)
				if len(sss[1]) == 0 {
					ff = cfg.ApplicationID + "-" + ff
				}
				log.Infof("including %s", ff)

				include = append(include, filepath.Join(targetSource, filepath.Base(ff)))
			}
			packageName := format.AddExt(fmt.Sprintf("%s-%s-from-%s", cfg.ApplicationID, ctx.Build.String(), from))
			targetPackage := fmt.Sprintf("target/%s", packageName)
			log.Infof("writing to '%s'", targetPackage)
			err = ar.Package(format, targetPackage, include)
			if err != nil {
				return errors.Wrapf(err, "packaging source files")
			}
			ctx.AddArtifact(packageName, targetPackage)
		}
	}

	log.Info("done")
	return nil
}

func GetPreviousTag(v string) (string, error) {
	commit, err := sh.Output("git", "rev-list", "--tags", "--skip=1", "--max-count=1")
	if err != nil {
		return "", errors.Wrapf(err, "getting last version commit")
	}
	fromTag, err := sh.Output("git", "describe", "--abbrev=0", commit)
	if err != nil {

		log.Warnf("No previous tag. Getting first commit")

		commit, err = sh.Output("git", "rev-list", "--max-parents=0", v)
		if err != nil {
			return "", errors.Wrapf(err, "getting first commit")
		}

		fromTag = commit
	}

	return fromTag, nil
}
