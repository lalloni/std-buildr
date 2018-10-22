package templates

import (
	"io/ioutil"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
)

func RenderEventualReadme(cfg *config.Config, location string) error {
	return ioutil.WriteFile(location, []byte(eveReadme), 0666)
}

const eveReadme = `# Eventuales
	
En esta carpeta se debe agregar los scripts eventuales segun los estandares.


**Ejemplo:**
	
    001-dml-descripción.sql
`
