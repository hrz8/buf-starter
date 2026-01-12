package views

import (
	"embed"
	"html/template"
	"io"
	"sync"
)

//go:embed *.html
var templateFS embed.FS

var (
	templates *template.Template
	once      sync.Once
	loadErr   error
)

func Load() error {
	once.Do(func() {
		templates, loadErr = template.ParseFS(templateFS, "*.html")
	})
	return loadErr
}

func Render(w io.Writer, name string, data any) error {
	if templates == nil {
		if err := Load(); err != nil {
			return err
		}
	}
	return templates.ExecuteTemplate(w, name, data)
}
