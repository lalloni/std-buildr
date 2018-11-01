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

func createEventualStructureOnBaseBranch(cfg *config.Config) error {
	cfg2 := *cfg
	cfg2.IssueID = ""
	return createEventualStructureOnBranch(&cfg2, localBaseBranch)
}

func createEventualStructureOnBranch(cfg *config.Config, branch string) error {
	if err := git.CreateOrphanBranch(branch); err != nil {
		return errors.Wrap(err, "creating base branch")
	}
	if err := createEventualStructure(cfg); err != nil {
		return nil
	}
	if err := git.CommitAddingAll("Crea estructura inicial para SQL eventual (std-buildr)"); err != nil {
		return errors.Wrap(err, "creating sql eventual commit in git")
	}
	return nil
}

func createEventualStructure(cfg *config.Config) error {
	if err := helpers.CreateProjectConfig(cfg); err != nil {
		return errors.Wrap(err, "creating project configuration")
	}
	if err := templates.RenderProjectReadme(cfg, templates.README); err != nil {
		return errors.Wrapf(err, "writing new project readme to %s", templates.README)
	}
	if err := createEventualReadme(cfg); err != nil {
		return errors.Wrap(err, "creating eventuals readme file")
	}
	return nil
}
