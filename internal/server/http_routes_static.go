package server

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/hrz8/altalune"
)

// registerStaticRoutes registers static file serving routes for the SPA frontend
func (s *Server) registerStaticRoutes(mux *http.ServeMux) {
	websiteFS, _ := fs.Sub(altalune.FrontendEmbeddedFiles, "frontend/.output/public")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Exclude API and OAuth endpoints from SPA serving
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/oauth/") {
			http.NotFound(w, r)
			return
		}

		p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")

		if p == "" {
			serveFileOr404(w, r, websiteFS, "index.html")
			return
		}

		if exists(websiteFS, p) {
			if isDir(websiteFS, p) {
				serveFileOr404(w, r, websiteFS, path.Join(p, "index.html"))
				return
			}
			serveFileOr404(w, r, websiteFS, p)
			return
		}

		serve404Page(w, r, websiteFS)
	})
}
