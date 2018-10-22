package templates

import (
	"io/ioutil"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
)

func RenderReplaceableReadme(cfg *config.Config, location string) error {
	return ioutil.WriteFile(location, []byte(repReadme), 0666)
}

const repReadme = `# Reemplazables

En esta carpeta se deben agregar los scripts reemplazables que seran incluidos en los scripts incrementales.
`
