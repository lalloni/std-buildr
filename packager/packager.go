package packager

import (
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/packager/sqleve"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/packager/sqlevo"
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
		return PackagerFunc(sqlevo.Package), nil
	case config.TypeOracleSQLEventual:
		return PackagerFunc(sqleve.Package), nil
	default:
		return nil, errors.Errorf("Packager not available for project type %q", cfg.Type)
	}
}
