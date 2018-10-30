package templates

import (
	"os"
	"text/template"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
)

func RenderProjectReadme(cfg *config.Config, location string) error {
	f, err := os.Open(location)
	if os.IsNotExist(err) {
		return render(projectReadmeTemplate, location, cfg)
	}
	if err == nil {
		f.Close()
	}
	return nil
}

var projectReadmeTemplate = template.Must(template.New("ProjectReadme").Parse(projectReadme))

const projectReadme = `# {{ .ApplicationID }}

Este proyecto esta siendo gestionado por [std-buildr](https://gitlab.cloudint.afip.gob.ar/std/std-buildr/).
`
