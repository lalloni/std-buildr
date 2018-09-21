package sqleve

import (
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/version"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/ar"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/sh"
)

const targetSource = "target/source"

func Package(cfg *config.Config, ctx *context.Context) error {

	ev,err := version.ParseEventualVersion(ctx.Build.Version)
	if err != nil {
		return errors.Wrap(err,"checking eventual version")
	}

	err = os.MkdirAll(targetSource, 0775)
	if err != nil {
		return errors.Wrapf(err, "creating directory")
	}
	base, err := sh.FirstExist("src/sql")
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Errorf("missing scripts sources (at 'src/sql')")
		}
		return errors.Wrapf(err, "checking eventual source presence")
	}

	// preprocess
	sourcesTargetMap := make(map[string]string)
	sources, err := sh.CollectFiles(base)
	if err != nil {
		return errors.Wrapf(err, "collecting source files from '%s'", base)
	}

	eventualMajorVersionRegexp := regexp.MustCompile(`^((.*)-([0-9]+))$`)

	v1 := eventualMajorVersionRegexp.FindStringSubmatch(ctx.Build.Version)

	version := fmt.Sprintf("^%s(-.*)?\\.sql$", v1[2])

	eventualRegexp := regexp.MustCompile(version)

	for _, source := range sources {

		if !eventualRegexp.MatchString(path.Base(source)) {
			return errors.Errorf("source file name '%s' does not match standard naming scheme (%s)", source, eventualRegexp.String())
		}

		targetName := cfg.ApplicationID + "-" + path.Base(source)
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
		in.Close()
		out.Close()
	}

	// package all
	log.Infof("generating full package")
	targetSources, err := sh.CollectFiles(targetSource)
	if err != nil {
		return errors.Wrapf(err, "collecting preprocessed files from '%s'", targetSource)
	}
	targetPackage := fmt.Sprintf("target/%s-%s.tar.xz", c.ApplicationID, ctx.Build.String())
	log.Infof("writing to '%s'", targetPackage)
	err = ar.TarXz(targetPackage, targetSources, path.Base)
	if err != nil {
		return errors.Wrapf(err, "packaging source files")
	}

	return nil
}
