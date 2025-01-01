package renderer

import (
	"errors"
	"html/template"
	"io"
)

var (
	errEmptyTemplateName = errors.New("empty template name")
)

// HTMLRenderer обслуживает рендер HTML шаблонов.
type HTMLRenderer struct {
	templates *template.Template
}

// NewHTMLRenderer создает рендерер.
func NewHTMLRenderer(templatePath string) *HTMLRenderer {
	return &HTMLRenderer{
		templates: template.Must(template.ParseFiles(templatePath)),
	}
}

// Render рендер шаблона с переданными данными.
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
