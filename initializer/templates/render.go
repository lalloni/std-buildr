package templates

import (
	"os"
	"text/template"
)

const README = "README.md"

func render(tpl *template.Template, target string, data interface{}) error {
	w, err := os.Create(target)
	if err != nil {
		return err
	}
	if err := tpl.Execute(w, data); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}
