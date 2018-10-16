package initializer

import (
	"fmt"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/filesContent"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/sqlevo"
)

type Initializer interface {
	Initialize(cfg *config.Config) error
}

type InitializerFunc func(cfg *config.Config) error

func (f InitializerFunc) Initialize(cfg *config.Config) error {

	cfg.SystemID = viper.GetString("buildr.system-id")

	if cfg.SystemID == "" {
		return errors.Errorf("system id is required")
	}

	cfg.ApplicationID = viper.GetString("buildr.application-id")
	if cfg.ApplicationID == "" {
		return errors.Errorf("application id is required")
	}

	if !strings.HasPrefix(cfg.ApplicationID, cfg.SystemID) {
		return errors.Errorf("system id (%s) must be prefix of application id (currently: %s)", cfg.SystemID, cfg.ApplicationID)
	}

	cfg.Type = viper.GetString("buildr.type")
	if cfg.Type == "" {
		return errors.Errorf("Type is required")
	}

	err := os.MkdirAll(cfg.ApplicationID, 0775)
	if err != nil {
		return errors.Wrapf(err, "creating directory %s", cfg.ApplicationID)
	}

	log.Infof("changing current directoty to %s", cfg.ApplicationID)
	err = os.Chdir("./" + cfg.ApplicationID)
	if err != nil {
		return errors.Wrapf(err, "changing current directory to %s", cfg.ApplicationID)
	}

	err = git.Initialize()
	if err != nil {
		return errors.Wrapf(err, "initializing git project")
	}

	buildrfile, err := yaml.Marshal(cfg)

	b, err := os.Create("buildr.yaml")
	if err != nil {
		return errors.Wrapf(err, "creating buildr.yaml")
	}
	defer b.Close()

	_, err = b.Write(buildrfile)
	if err != nil {
		return errors.Wrapf(err, "writing buildr.yaml")
	}

	rootReadme, err := os.Create("readme.md")
	if err != nil {
		return errors.Wrapf(err, "creating readme.md")
	}
	defer rootReadme.Close()
	content := fmt.Sprintf(filesContent.README_ROOT, cfg.ApplicationID)
	_, err = rootReadme.Write([]byte(content))
	if err != nil {
		return errors.Wrapf(err, "writing %s", "readme.md")
	}

	remote := fmt.Sprintf("git@gitlab.cloudint.afip.gob.ar:%s/%s.git", cfg.SystemID, cfg.ApplicationID)

	err = git.AddRemote(remote)
	if err != nil {
		return errors.Wrapf(err, "adding git remote (%s)", remote)
	}

	fs, err := git.AddAll()
	if err != nil {
		return errors.Wrapf(err, "adding files (%s)", fs)
	}

	return f(cfg)
}

func New(cfg *config.Config) (Initializer, error) {
	cfg.Type = viper.GetString("buildr.type")
	switch cfg.Type {
	case config.TypeOracleSQLEvolutional:
		return InitializerFunc(sqlevo.Initialize), nil
	default:
		return nil, errors.Errorf("Initializer not available for project type %q", cfg.Type)
	}
}
