package publisher

import (
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/context"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/publisher/sql"
)

type Publisher interface {
	Publish(cfg *config.Config, ctx *context.Context) error
}

type PublisherFunc func(cfg *config.Config, ctx *context.Context) error

func (f PublisherFunc) Publish(cfg *config.Config, ctx *context.Context) error {
	return f(cfg, ctx)
}

func New(cfg *config.Config) (Publisher, error) {
	switch cfg.Type {
	case config.TypeOracleSQLEvolutional, config.TypeOracleSQLDeferred, config.TypeOracleSQLEventual:
		return PublisherFunc(sql.Publish), nil
	default:
		return nil, errors.Errorf("Publisher not available for project type %q", cfg.Type)
	}
}
