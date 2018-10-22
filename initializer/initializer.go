package initializer

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/sqldef"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/sqleve"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/sqlevo"
)

type Initializer interface {
	Initialize(cfg *config.Config) error
}

type InitializerFunc func(cfg *config.Config) error

func (f InitializerFunc) Initialize(cfg *config.Config) error {
	return f(cfg)
}

func New(cfg *config.Config) (Initializer, error) {
	cfg.Type = viper.GetString("buildr.type")
	switch cfg.Type {
	case config.TypeOracleSQLEvolutional:
		return InitializerFunc(sqlevo.Initialize), nil
	case config.TypeOracleSQLDeferred:
		return InitializerFunc(sqldef.Initialize), nil
	case config.TypeOracleSQLEventual:
		return InitializerFunc(sqleve.Initialize), nil
	default:
		return nil, errors.Errorf("Initializer not available for project type '%s'", cfg.Type)
	}
}

func CreateEventual(cfg *config.Config) error {

	cfg.IssueID = viper.GetString("buildr.issue-id")

	if cfg.IssueID == "" {
		return errors.Errorf("issue id is required")
	}

	return sqleve.CreateEventual(cfg)

}
