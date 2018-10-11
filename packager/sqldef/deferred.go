package sqldef

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/ar"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/sh"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/version"
)

var (
	tagNameRegexp  = regexp.MustCompile(`^v(.*)$`)
	deferredRegexp = regexp.MustCompile(`^(.*)\.sql$`)
)

func Package(cfg *config.Config, ctx *context.Context) error {

	git.GetStateIn(ctx)

	v := tagNameRegexp.FindStringSubmatch(ctx.Build.Version)
	if v == nil {
		return errors.Errorf("tag name must be prefixed with a 'v' character (found '%s')", ctx.Build.Version)
	}
	ctx.Build.Version = v[1]

	_, err := version.ParseSemanticVersion(ctx.Build.Version)
	if err != nil {
		return errors.Wrap(err, "checking valid semantic version")
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

		if !deferredRegexp.MatchString(path.Base(source)) {
			return errors.Errorf("source file name '%s' does not match standard naming scheme (%s)", source, deferredRegexp.String())
		}

		targetName := path.Base(source)

		exp := fmt.Sprintf("^%s-(.*)$", cfg.ApplicationID)
		appDeferredRegexp := regexp.MustCompile(exp)

		if !appDeferredRegexp.MatchString(path.Base(source)) {
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
	packageName := fmt.Sprintf("%s-%s.tar.xz", cfg.ApplicationID, ctx.Build.String())
	targetPackage := fmt.Sprintf("target/%s", packageName)
	log.Infof("writing to '%s'", targetPackage)
	err = ar.TarXz(targetPackage, targetSources, path.Base)
	if err != nil {
		return errors.Wrapf(err, "packaging source files")
	}
	ctx.AddArtifact(packageName, targetPackage)

	return nil
}
