package packager

import (
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/packages"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
)

type Packager interface {
	Package(cfg *config.Config, ctx *context.Context) error
}

type PackagerFunc func(cfg *config.Config, ctx *context.Context) error

func (f PackagerFunc) Package(cfg *config.Config, ctx *context.Context) error {
	return f(cfg, ctx)
}

func New(cfg *config.Config) (Packager, error) {
	switch cfg.Type {
	case config.TypeOracleSQLEvolutional:
		return PackagerFunc(packages.PackageOracleSQLEvolutional), nil
	case config.TypeOracleSQLEventual:
		return PackagerFunc(packages.PackageOracleSQLEventual), nil
	default:
		return nil, errors.Errorf("Packager not available for project type %q", cfg.Type)
	}
}
