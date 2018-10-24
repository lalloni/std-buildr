package templates

import (
	"io/ioutil"
	"os"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
)

func RenderEventualReadme(cfg *config.Config, location string) error {
	f, err := os.Open(location)
	if os.IsNotExist(err) {
		return ioutil.WriteFile(location, []byte(eveReadme), 0666)
	}
	if err == nil {
		f.Close()
	}
	return nil
}

const eveReadme = `# Eventuales
	
En esta carpeta se debe agregar los scripts eventuales segun los estandares.


**Ejemplo:**
	
    001-dml-descripción.sql
`
