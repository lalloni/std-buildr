package templates

import (
	"text/template"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/config"
)

func RenderProjectReadme(cfg *config.Config, location string) error {
	return render(projectReadmeTemplate, location, cfg)
}

var projectReadmeTemplate = template.Must(template.New("ProjectReadme").Parse(projectReadme))

const projectReadme = `# {{ .ApplicationID }}

Este proyecto esta siendo gestionado por [std-buildr](https://gitlab.cloudint.afip.gob.ar/std/std-buildr/blob/master/doc/index.md).
`
