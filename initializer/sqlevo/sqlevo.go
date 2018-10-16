package sqlevo

import (
	"os"

	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/filesContent"
	"gitlab.cloudint.afip.gob.ar/std/std-buildr/git"
)

const (
	incFolder = "src/sql/inc"
	repFolder = "src/sql/rep"
)

func Initialize(cfg *config.Config) error {

	err := os.MkdirAll(incFolder, 0775)
	if err != nil {
		return errors.Wrapf(err, "creating directory")
	}

	incReadme, err := os.Create(incFolder + "/readme.md")
	if err != nil {
		return errors.Wrapf(err, "creating %s", incFolder+"/readme.md")
	}
	defer incReadme.Close()

	_, err = incReadme.Write([]byte(filesContent.README_INC))
	if err != nil {
		return errors.Wrapf(err, "writing %s", incFolder+"/readme.md")
	}

	err = os.MkdirAll(repFolder, 0775)
	if err != nil {
		return errors.Wrapf(err, "creating directory")
	}

	repReadme, err := os.Create(repFolder + "/readme.md")
	if err != nil {
		return errors.Wrapf(err, "creating %s", repFolder+"/readme.md")
	}
	defer repReadme.Close()

	_, err = repReadme.Write([]byte(filesContent.README_REP))
	if err != nil {
		return errors.Wrapf(err, "writing %s", repFolder+"/readme.md")
	}

	uu, err := git.AddAll()
	if err != nil {
		return errors.Wrapf(err, "adding  %s", uu)
	}

	err = git.Commit("First Commit by std-buildr")
	if err != nil {
		return errors.Wrapf(err, "commiting")
	}

	return nil
}
