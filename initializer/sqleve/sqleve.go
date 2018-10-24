package sqleve

import (
	"fmt"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/initializer/helpers"

	"github.com/apex/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
)

const (
	eveFolder        = "src/sql"
	localBaseBranch  = "base"
	remoteBaseBranch = "origin/base"
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
		return errors.New("tracker-id is required")
	}
	cfg.TrackerID = tid

	if err := helpers.CreateProject(cfg); err != nil {
		return errors.Wrap(err, "creating project structure")
	}

	if err := git.CommitAddingAll("Crea estructura inicial de proyecto SQL eventual (std-buildr)"); err != nil {
		return errors.Wrap(err, "creating sql eventual commit in git")
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

	newBranch := fmt.Sprintf("%s-%s", cfg.TrackerID, cfg.IssueID)

	// validar que exista branch base
	existLocal, err := git.ExistBranch(localBaseBranch)
	if err != nil {
		return errors.Wrap(err, "checking for local base branch existence")
	}
	if existLocal {
		log.Warnf("new eventual branch will be based on the local branch named '%s'", localBaseBranch)
	} else {
		existRemote, err := git.ExistBranch(remoteBaseBranch)
		if err != nil {
			return errors.Wrap(err, "checking for remote base branch existence")
		}
		if existRemote {
			log.Warnf("new eventual branch will be based on the remote branch named '%s'", remoteBaseBranch)
			// crear local desde remote
			err := git.CreateBranchFrom(localBaseBranch, remoteBaseBranch)
			if err != nil {
				return errors.Wrap(err, "creating local base branch from remote base branch")
			}
			existLocal = true
		} else {
			// solo creamos el branch del eventual
			if err := git.CreateOrphanBranch(newBranch); err != nil {
				return errors.Wrap(err, "creating base branch")
			}
		}
	}

	// si existía o creamos la local base...
	if existLocal {
		// crear branch para el eventual desde base local
		if err := git.CreateBranchFrom(newBranch, localBaseBranch); err != nil {
			return errors.Wrapf(err, "creating branch %s from %s", newBranch, localBaseBranch)
		}
	}

	// en este punto debe existir branch del eventual y solo tratamos con ella en adelante

	// creamos estructura de archivos para el eventual
	if err := createEventualStructure(cfg); err != nil {
		return errors.Wrap(err, "creating eventual structure")
	}

	// commiteamos estado en el branch del eventual
	if err := git.CommitAddingAll("Crea estructura inicial para SQL eventual (std-buildr)"); err != nil {
		return errors.Wrap(err, "creating sql eventual commit in git")
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
