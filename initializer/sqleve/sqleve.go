package sqleve

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/templates"
)

const (
	eveFolder = "src/sql"
)

func Initialize(cfg *config.Config) error {

	cfg.TrackerID = viper.GetString("buildr.tracker-id")

	if cfg.TrackerID == "" {
		return errors.Errorf("--tracker-id is required")
	}

	err := os.MkdirAll(eveFolder, 0775)
	if err != nil {
		return errors.Wrapf(err, "creating directory")
	}

	p := filepath.Join(eveFolder, templates.README)
	if err := templates.RenderEventualReadme(cfg, p); err != nil {
		return errors.Wrapf(err, "creating %s", p)
	}

	git.CreateBranch("base")
	if err != nil {
		return errors.Wrapf(err, "creating branch %s", "base")
	}

	buildrfile, err := yaml.Marshal(cfg)

	b, err := os.Create("buildr.yaml")
	if err != nil {
		return errors.Wrapf(err, "opening buildr.yaml")
	}
	defer b.Close()

	_, err = b.Write(buildrfile)
	if err != nil {
		return errors.Wrapf(err, "writing buildr.yaml")
	}

	if err := git.CommitAddingAll("Crea estructura de proyecto eventual (std-buildr)"); err != nil {
		return errors.Wrap(err, "creating sql eventual commit in git")
	}

	git.Push("origin", "base")

	return nil
}
