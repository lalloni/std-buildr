package packages

import (
	"fmt"
	"path"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/apex/log"
	"github.com/pkg/errors"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/ar"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/sh"
)

func VerifyEventualSQLOracleVersion(ctx *context.Context) error {

	versionRegex := regexp.MustCompile(`^((.*)-([0-9]+)-([0-9]+))(-([0-9]+)-(.*))?$`)

	v1 := versionRegex.FindStringSubmatch(ctx.Build.Version)

	if v1 == nil {
		return errors.Errorf("invalid tag name")
	}

	if len(v1[5]) > 0 {
		ctx.Build.Prerelease = v1[5]
	}
	ctx.Build.Version = v1[1]

	log.Infof("version is '%s'", ctx.Build.Version)

	return nil

}

func VerifyStandardVersion(ctx *context.Context) error {

	v2 := tagNameRegexp.FindStringSubmatch(ctx.Build.Version)
	if v2 == nil {
		return errors.Errorf("tag name must be prefixed with a 'v' character (found '%s')", ctx.Build.Version)
	}
	version, err := semver.NewVersion(v2[1])
	if err != nil {
		return errors.Wrapf(err, "tag name must be a valid semver 2 string prefixed with a 'v' character (found '%s')", ctx.Build.Version)
	}

	ctx.Build.Version = fmt.Sprintf("%d.%d.%d", version.Major(), version.Minor(), version.Patch())
	ctx.Build.Prerelease = version.Prerelease()
	log.Infof("version is '%s'", version)
	return nil

}

type packageSettings struct {
	Context      *context.Context
	Config       *config.Config
	TargetSource string
}

func PackageAllSQL(targetSource string, ctx *context.Context, c *config.Config) error {

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
	ctx.AddArtifact(targetPackage)

	return nil
}
