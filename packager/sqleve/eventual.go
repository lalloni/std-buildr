package sqleve

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

const targetSource = "target/source"

var eventualRegexp = regexp.MustCompile(`^(?:(.*)([-_]))?(\d{3,})([-_])(dml|dcl|ddl)(?:([-_])(.+))?\.sql$`)

func Package(cfg *config.Config, ctx *context.Context) error {

	if cfg.IssueID == "" {
		return errors.Errorf(`issue-id is required in configuration`)
	}

	ev, err := version.ParseEventualVersion(ctx.Build.Version)
	if err != nil {
		return errors.Wrap(err, "checking eventual version")
	}

	tagVersion := fmt.Sprintf("%s-%s-%d", ev.TrackerID, ev.IssueID, ev.Version)
	configVersion := fmt.Sprintf("%s-%s-%d", cfg.TrackerID, cfg.IssueID, ev.Version)
	if ev.TrackerID != cfg.TrackerID {

		return errors.Errorf(msg.PACKAGE_EVENTUAL_INVALID_TRACKER, ev.TrackerID, tagVersion, cfg.TrackerID, configVersion, configVersion, tagVersion, tagVersion)
	}
	if ev.IssueID != cfg.IssueID {

		return errors.Errorf(msg.PACKAGE_EVENTUAL_INVALID_ISSUE, ev.IssueID, tagVersion, cfg.IssueID, configVersion, configVersion, tagVersion, tagVersion)
	}

	allowDirty := viper.GetBool("buildr.allow-dirty")
	if ctx.Build.Dirty() && !allowDirty {
		untrackedAndChangedFiles, err := git.ListUntrackedFilesAndChangedFiles()
		if err != nil {
			return errors.Wrap(err, "checking untracked and changed files")
		}
		uu := strings.Join(untrackedAndChangedFiles, " ")
		if uu == "" {
			return errors.Errorf(msg.PACKAGE_EVENTUAL_COMMITED_ERROR, tagVersion, tagVersion, tagVersion)
		}
		return errors.Errorf(msg.PACKAGE_EVENTUAL_UNTRACKER_ERROR, uu, cfg.TrackerID, cfg.IssueID, tagVersion, tagVersion, tagVersion)
	}

	allowUntagged := viper.GetBool("buildr.allow-untagged")
	if ctx.Build.Untagged() && !allowUntagged {
		return errors.Errorf(msg.PACKAGE_EVENTUAL_UNTAGGED_ERROR, cfg.TrackerID, cfg.IssueID, cfg.TrackerID, cfg.IssueID, tagVersion, tagVersion, tagVersion)
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

	for _, source := range sources {

		if filepath.Base(source) == templates.README {
			continue
		}
		if !eventualRegexp.MatchString(filepath.Base(source)) {
			return errors.Errorf("source file name '%s' does not match standard naming scheme (%s)", source, eventualRegexp.String())
		}

		ss := eventualRegexp.FindStringSubmatch(filepath.Base(source))
		tag := fmt.Sprintf("%s-%s-%d", ev.TrackerID, ev.IssueID, ev.Version)
		if len(ss[1]) != 0 && (ss[1] != tag) {
			tag2 := fmt.Sprintf("%s_%s_%d", ev.TrackerID, ev.IssueID, ev.Version)
			if ss[1] != tag2 {
				return errors.Errorf("source file '%s' name prefix '%s' must equal tag '%s' if used", source, ss[1], tag)
			}
			fn := fmt.Sprintf("%s-%s-%s-%s.sql", tag, ss[3], ss[5], ss[7])
			log.Warnf(msg.PACKAGE_EVENTUAL_FILENAME_WARN, ss[0], fn)
		} else {
			if ss[2] == "_" || ss[4] == "_" || ss[6] == "_" {
				fn := fmt.Sprintf("%s-%s-%s-%s.sql", tag, ss[3], ss[5], ss[7])
				log.Warnf(msg.PACKAGE_EVENTUAL_FILENAME_WARN, ss[0], fn)
			}

		}

		targetName := filepath.Base(source)
		if len(ss[1]) == 0 {
			targetName = tag + "-" + targetName
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
