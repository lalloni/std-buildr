package sqleve

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/helpers"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/templates"
)

func createEventualReadme(cfg *config.Config) error {
	err := os.MkdirAll(eveFolder, 0775)
	if err != nil {
		return errors.Wrapf(err, "creating directory")
	}
	p := filepath.Join(eveFolder, templates.README)
	if err := templates.RenderEventualReadme(cfg, p); err != nil {
		return errors.Wrapf(err, "creating %s", p)
	}
	return nil
}

func createBaseBranch(cfg *config.Config) error {
	if err := git.CreateOrphanBranch(baseBranch); err != nil {
		return errors.Wrap(err, "creating base branch")
	}
	cfg2 := *cfg
	cfg2.IssueID = ""
	if err := helpers.CreateProjectConfig(&cfg2); err != nil {
		return errors.Wrap(err, "creating project configuration")
	}
	if err := templates.RenderProjectReadme(cfg, templates.README); err != nil {
		return errors.Wrapf(err, "writing new project readme to %s", templates.README)
	}
	if err := createEventualReadme(cfg); err != nil {
		return errors.Wrap(err, "creating eventuals readme file")
	}
	if err := git.CommitAddingAll("Crea estructura inicial para SQL eventual (std-buildr)"); err != nil {
		return errors.Wrap(err, "creating sql eventual commit in git")
	}
	if err := git.Push("origin", "base"); err != nil {
		return errors.Wrap(err, "pushing base branch to origin")
	}
	return nil
}
