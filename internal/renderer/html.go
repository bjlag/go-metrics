package renderer

import (
	"errors"
	"html/template"
	"io"
)

var (
	errEmptyTemplateName = errors.New("empty template name")
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
		return errEmptyTemplateName
	}

	err := r.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		return err
	}

	return nil
}
