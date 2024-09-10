package renderer

import (
	"errors"
	"html/template"
	"io"
)

type HTMLRenderer struct {
	templates *template.Template
}

func NewHTMLRenderer(templatePath string) *HTMLRenderer {
	return &HTMLRenderer{
		templates: template.Must(template.ParseFiles(templatePath)),
	}
}

func (r HTMLRenderer) Render(w io.Writer, name string, data interface{}) error {
	if len(name) == 0 {
		return errors.New("empty name template")
	}

	err := r.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		return err
	}

	return nil
}
