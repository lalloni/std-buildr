package sqleve

import (
	"fmt"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/helpers"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
)

const (
	eveFolder  = "src/sql"
	baseBranch = "base"
)

type Script struct {
	Name string
	Type string
}

func Initialize(cfg *config.Config) error {

	if err := helpers.ValidateProjectConfig(cfg); err != nil {
		return errors.Wrap(err, "checking project configuration requirements")
	}

	tid := viper.GetString("buildr.tracker-id")
	if tid == "" {
		return errors.Errorf("--tracker-id is required")
	}
	cfg.TrackerID = tid

	if err := helpers.CreateProject(cfg); err != nil {
		return errors.Wrap(err, "creating project structure")
	}

	if err := createBaseBranch(cfg); err != nil {
		return errors.Wrap(err, "creating eventual base branch")
	}

	return nil
}

func CreateEventual(cfg *config.Config) error {

	// validar que no haya cambios en la WC
	untracked, err := git.UntrackedFilesInCWD()
	if err != nil {
		return errors.Wrapf(err, "checking for untracked files")
	}
	if untracked || git.ChangedFilesInCWD() || git.UncommittedFilesInCWD() {
		return errors.Errorf("you have changes uncommited, please commit them or undo")
	}

	// validar que exista branch base
	exist, err := git.ExistBranch(baseBranch)
	if err != nil {
		return errors.Wrap(err, "checking for base branch existence")
	}
	if !exist {
		exist, err := git.ExistBranch(baseBranch)
		if err != nil {
			return errors.Wrap(err, "checking for remote base branch existence")
		}
		if exist {
			git.CreateBranchFrom("base", "origin/base")
		} else {
			if err := createBaseBranch(cfg); err != nil {
				return errors.Wrap(err, "creating eventual base branch")
			}
		}
	}

	// crear branch para el eventual
	newBranch := fmt.Sprintf("%s-%s", cfg.TrackerID, cfg.IssueID)
	if err := git.CreateBranchFrom(newBranch, baseBranch); err != nil {
		return errors.Wrapf(err, "creating branch %s from %s", newBranch, baseBranch)
	}

	// crear config de proyecto
	if err := helpers.CreateProjectConfig(cfg); err != nil {
		return errors.Wrap(err, "creating project configuration")
	}

	// crear scripts solicitados
	ss := viper.Get("buildr.scripts").([]Script)
	for i, script := range ss {
		err := helpers.CreateEmptyFilef("%s/%03d-%s-%s.sql", eveFolder, i+1, script.Type, script.Name)
		if err != nil {
			return errors.Wrapf(err, "creating file")
		}
	}

	return nil
}
