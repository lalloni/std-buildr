package sqleve

import (
	"fmt"
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

type Script struct {
	Name string
	Type string
}

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

func CreateEventual(cfg *config.Config) error {

	baseBranch := "base"
	untracked, err := git.UntrackedFilesInCWD()
	if err != nil {

		return errors.Wrapf(err, "checking for untracked files")
	}

	if untracked || git.ChangedFilesInCWD() || git.UncommittedFilesInCWD() {
		return errors.Errorf("you have changes uncommited, please commit them or undo")
	}

	if !git.CheckExistingBranch(baseBranch) {
		if !git.CheckExistingBranch("origin/base") {

			err := os.RemoveAll(eveFolder)
			if err != nil {
				return errors.Wrapf(err, "removing %s", eveFolder)
			}
			err = Initialize(cfg)
			if err != nil {
				return errors.Wrapf(err, "initializing project")
			}
		} else {
			baseBranch = "origin/base"
		}
	} else {

		buildrfile, err := yaml.Marshal(cfg)

		b, err := os.OpenFile("buildr.yaml", os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return errors.Wrapf(err, "opening buildr.yaml")
		}
		defer b.Close()

		_, err = b.Write(buildrfile)
		if err != nil {
			return errors.Wrapf(err, "writing buildr.yaml")
		}

	}

	ss := viper.Get("buildr.scripts").([]Script)
	for i := 0; i < len(ss); i++ {
		file := fmt.Sprintf("%s/%03d-%s-%s.sql", eveFolder, i+1, ss[i].Type, ss[i].Name)
		createFile(file)
	}
	newBranch := fmt.Sprintf("%s-%s", cfg.TrackerID, cfg.IssueID)
	err = git.CreateBranchFrom(newBranch, baseBranch)
	if err != nil {
		return errors.Wrapf(err, "creaating branch %s from %s", newBranch, baseBranch)
	}
	fs, err := git.AddAll()
	if err != nil {
		return errors.Wrapf(err, "adding untracked and changed files %s", fs)
	}
	return nil
}

func createFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return errors.Wrapf(err, "creating %s", file)
	}
	defer f.Close()
	_, err = f.Write([]byte(""))
	if err != nil {
		return errors.Wrapf(err, "writing %s", file)
	}
	return nil
}
