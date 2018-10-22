package sqlevo

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/templates"
)

var (
	incFolder = filepath.Join("src", "sql", "inc")
	repFolder = filepath.Join("src", "sql", "rep")
)

func Initialize(cfg *config.Config) error {

	if err := os.MkdirAll(incFolder, 0775); err != nil {
		return errors.Wrap(err, "creating incremental sources directory")
	}

	p := filepath.Join(incFolder, templates.README)
	if err := templates.RenderIncrementalReadme(cfg, p); err != nil {
		return errors.Wrapf(err, "creating %s", p)
	}

	if err := os.MkdirAll(repFolder, 0775); err != nil {
		return errors.Wrap(err, "creating replaceable sources directory")
	}

	p = filepath.Join(repFolder, templates.README)
	if err := templates.RenderReplaceableReadme(cfg, p); err != nil {
		return errors.Wrapf(err, "creating %s", p)
	}

	if err := git.CommitAddingAll("Crea estructura de proyecto evolutivo (std-buildr)"); err != nil {
		return errors.Wrap(err, "creating sql evolutional commit in git")
	}

	return nil

}
