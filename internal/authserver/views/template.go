package views

import (
	"embed"
	"html/template"
	"io"
	"strings"
	"sync"
	"time"
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
		funcMap := template.FuncMap{
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
			"substr": func(s string, start, length int) string {
				if start >= len(s) {
					return ""
				}
				end := start + length
				if end > len(s) {
					end = len(s)
				}
				return strings.ToUpper(s[start:end])
			},
			"split": func(s, sep string) []string {
				if s == "" {
					return []string{}
				}
				return strings.Split(s, sep)
			},
			"formatTime": func(t time.Time) string {
				return t.Format("Jan 2, 2006 at 3:04 PM")
			},
		}
		templates, loadErr = template.New("").Funcs(funcMap).ParseFS(templateFS, "*.html")
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
