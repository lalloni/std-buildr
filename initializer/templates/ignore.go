package templates

import (
	"io/ioutil"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
)

func RenderGitIgnore(cfg *config.Config, location string) error {
	return ioutil.WriteFile(location, []byte(gitignore), 0666)
}

const gitignore = `target/`
