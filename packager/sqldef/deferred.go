package sqldef

import (
	"bufio"
	"fmt"
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
	tagNameRegexp  = regexp.MustCompile(`^v(.*)$`)
	deferredRegexp = regexp.MustCompile(`^(.*)\.sql$`)
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
	base, err := sh.FirstExist("src/sql")
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Errorf("missing scripts sources (at 'src/sql')")
		}
		return errors.Wrapf(err, "checking scripts source presence")
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

		if !deferredRegexp.MatchString(filepath.Base(source)) {
			return errors.Errorf("source file name '%s' does not match standard naming scheme (%s)", source, deferredRegexp.String())
		}

		targetName := filepath.Base(source)

		if !strings.HasPrefix(targetName, cfg.ApplicationID) {
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
			_, err := fmt.Fprintln(out, l)
			if err != nil {
				return errors.Wrap(err, "copying input to output")
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

	return nil
}
