package templates

import (
	"io/ioutil"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
)

func RenderIncrementalReadme(cfg *config.Config, location string) error {
	return ioutil.WriteFile(location, []byte(incReadme), 0666)
}

const incReadme = `# Incrementales

En esta carpeta se deben agregar los scripts incrementales segun los estandares.

**Ejemplo:**

    000001-dml-descripción.sql
`
