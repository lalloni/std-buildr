package sqldef

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/helpers"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/templates"
)

const srcFolder = "src/sql"

func Initialize(cfg *config.Config) error {

	if err := helpers.ValidateProjectConfig(cfg); err != nil {
		return errors.Wrap(err, "checking project configuration requirements")
	}

	if err := helpers.CreateProject(cfg, "evolutivo"); err != nil {
		return errors.Wrap(err, "creating project structure")
	}

	if err := os.MkdirAll(srcFolder, 0775); err != nil {
		return errors.Wrapf(err, "creating directory")
	}

	p := filepath.Join(srcFolder, templates.README)
	if err := templates.RenderDeferredReadme(cfg, p); err != nil {
		return errors.Wrapf(err, "creating %s", p)
	}

	if err := git.CommitAddingAll("Crea estructura inicial de SQL diferido (std-buildr)"); err != nil {
		return errors.Wrap(err, "creating sql deferred commit in git")
	}

	return nil
}
