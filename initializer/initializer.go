package initializer

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/sqldef"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/sqlevo"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/templates"
)

type Initializer interface {
	Initialize(cfg *config.Config) error
}

type InitializerFunc func(cfg *config.Config) error

func (f InitializerFunc) Initialize(cfg *config.Config) error {

	cfg.SystemID = viper.GetString("buildr.system-id")

	if cfg.SystemID == "" {
		return errors.New("system id is required")
	}

	cfg.ApplicationID = viper.GetString("buildr.application-id")
	if cfg.ApplicationID == "" {
		return errors.New("application id is required")
	}

	if !strings.HasPrefix(cfg.ApplicationID, cfg.SystemID) {
		return errors.Errorf("system id (%s) must be prefix of application id (currently: %s)", cfg.SystemID, cfg.ApplicationID)
	}

	cfg.Type = viper.GetString("buildr.type")
	if cfg.Type == "" {
		return errors.New("project type is required")
	}

	if err := os.MkdirAll(cfg.ApplicationID, 0775); err != nil {
		return errors.Wrapf(err, "creating directory '%s'", cfg.ApplicationID)
	}

	if err := os.Chdir(cfg.ApplicationID); err != nil {
		return errors.Wrapf(err, "changing current directory to '%s'", cfg.ApplicationID)
	}

	if err := git.Init(); err != nil {
		return errors.Wrap(err, "initializing git repository")
	}

	config, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "marshalling new project configuration")
	}

	if err := ioutil.WriteFile("buildr.yaml", config, 0666); err != nil {
		return errors.Wrap(err, "writing new project configuration to buildr.yaml")
	}

	if err := templates.RenderProjectReadme(cfg, templates.README); err != nil {
		return errors.Wrapf(err, "writing new project readme to %s", templates.README)
	}

	remote := fmt.Sprintf("git@gitlab.cloudint.afip.gob.ar:%s/%s.git", cfg.SystemID, cfg.ApplicationID)
	if err := git.AddRemote("origin", remote); err != nil {
		return errors.Wrapf(err, "adding git remote (%s)", remote)
	}

	if err = git.CommitAddingAll("Crea estructura inicial (std-buildr)"); err != nil {
		return errors.Wrap(err, "creating initial commit in git")
	}

	return f(cfg)

}

func New(cfg *config.Config) (Initializer, error) {
	cfg.Type = viper.GetString("buildr.type")
	switch cfg.Type {
	case config.TypeOracleSQLEvolutional:
		return InitializerFunc(sqlevo.Initialize), nil
	case config.TypeOracleSQLDeferred:
		return InitializerFunc(sqldef.Initialize), nil
	default:
		return nil, errors.Errorf("Initializer not available for project type '%s'", cfg.Type)
	}
}
