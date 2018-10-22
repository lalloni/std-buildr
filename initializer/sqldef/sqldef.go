package sqldef

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/templates"
)

const srcFolder = "src/sql"

func Initialize(cfg *config.Config) error {

	if err := os.MkdirAll(srcFolder, 0775); err != nil {
		return errors.Wrapf(err, "creating directory")
	}

	p := filepath.Join(srcFolder, templates.README)
	if err := templates.RenderDeferredReadme(cfg, p); err != nil {
		return errors.Wrapf(err, "creating %s", p)
	}

	if err := git.CommitAddingAll("Crea estructura de proyecto diferido (std-buildr)"); err != nil {
		return errors.Wrap(err, "creating sql deferred commit in git")
	}

	return nil
}
