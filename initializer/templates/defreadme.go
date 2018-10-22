package templates

import (
	"io/ioutil"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
)

func RenderDeferredReadme(cfg *config.Config, location string) error {
	return ioutil.WriteFile(location, []byte(defReadme), 0666)
}

const defReadme = `# Diferidos

En esta carpeta se debe agregar los scripts para la aplicacion SQL Oracle Diferida.
`
