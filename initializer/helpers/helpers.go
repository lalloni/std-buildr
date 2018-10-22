package helpers

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/templates"
)

func CreateProject(cfg *config.Config) error {

	if err := os.MkdirAll(cfg.ApplicationID, 0775); err != nil {
		return errors.Wrapf(err, "creating directory '%s'", cfg.ApplicationID)
	}

	if err := os.Chdir(cfg.ApplicationID); err != nil {
		return errors.Wrapf(err, "changing current directory to '%s'", cfg.ApplicationID)
	}

	if err := git.Init(); err != nil {
		return errors.Wrap(err, "initializing git repository")
	}

	if err := CreateProjectConfig(cfg); err != nil {
		return errors.Wrap(err, "creating project configuration")
	}

	if err := templates.RenderProjectReadme(cfg, templates.README); err != nil {
		return errors.Wrapf(err, "writing new project readme to %s", templates.README)
	}

	remote := fmt.Sprintf("git@gitlab.cloudint.afip.gob.ar:%s/%s.git", cfg.SystemID, cfg.ApplicationID)
	if err := git.AddRemote("origin", remote); err != nil {
		return errors.Wrapf(err, "adding git remote (%s)", remote)
	}

	return nil

}

func ValidateProjectConfig(cfg *config.Config) error {

	systemID := viper.GetString("buildr.system-id")
	if systemID == "" {
		return errors.New("system id is required")
	}

	applicationID := viper.GetString("buildr.application-id")
	if applicationID == "" {
		return errors.New("application id is required")
	}

	if !strings.HasPrefix(applicationID, systemID) {
		return errors.Errorf("system id (%s) must be prefix of application id (currently: %s)", systemID, applicationID)
	}

	projectType := viper.GetString("buildr.type")
	if projectType == "" {
		return errors.New("project type is required")
	}

	cfg.SystemID = systemID
	cfg.ApplicationID = applicationID
	cfg.Type = projectType

	return nil
}

func CreateEmptyFilef(format string, args ...interface{}) error {
	var (
		f   io.Closer
		err error
	)
	file := fmt.Sprintf(format, args...)
	if f, err = os.Create(file); err == nil {
		err = f.Close()
	}
	return err
}

func CreateProjectConfig(cfg *config.Config) error {
	var (
		bs  []byte
		err error
	)
	if bs, err = yaml.Marshal(cfg); err == nil {
		err = ioutil.WriteFile("buildr.yaml", bs, 0666)
	}
	return err
}
