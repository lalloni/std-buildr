package sqleve

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/version"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/ar"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/sh"
)

const targetSource = "target/source"

func Package(cfg *config.Config, ctx *context.Context) error {

	ev, err := version.ParseEventualVersion(ctx.Build.Version)
	if err != nil {
		return errors.Wrap(err, "checking eventual version")
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

	scriptName := fmt.Sprintf("^(.*-)?%s-%s-(dml|dcl|ddl)(-.*)?\\.sql$", ev.TrackerID, ev.IssueID)

	scriptNameRegexp := regexp.MustCompile(scriptName)

	for _, source := range sources {

		if !scriptNameRegexp.MatchString(path.Base(source)) {
			return errors.Errorf("source file name '%s' does not match standard naming scheme (%s)", source, scriptNameRegexp.String())
		}

		ss := scriptNameRegexp.FindStringSubmatch(path.Base(source))
		if len(ss[1]) != 0 && ss[1] != cfg.ApplicationID+"-" {
			return errors.Errorf("source file '%s' name prefix '%s' must equal application id '%s' if used", source, ss[1][:len(ss[1])-1], cfg.ApplicationID)
		}

		targetName := path.Base(source)
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
			_, err := fmt.Fprintln(out, l)
			if err != nil {
				return errors.Wrap(err, "copying input to output")
			}
		}

		in.Close()
		out.Close()
	}

	// package all
	log.Infof("generating full package")
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

	ctx.AddArtifact(packageName, targetPackage, (ev.Prerelease != "" || ctx.Build.Dirty()))

	return nil
}
