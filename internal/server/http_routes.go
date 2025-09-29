package server

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/gen/altalune/v1/altalunev1connect"
	"github.com/hrz8/altalune/gen/greeter/v1/greeterv1connect"
	api_key_domain "github.com/hrz8/altalune/internal/domain/api_key"
	employee_domain "github.com/hrz8/altalune/internal/domain/employee"
	greeter_domain "github.com/hrz8/altalune/internal/domain/greeter"
	project_domain "github.com/hrz8/altalune/internal/domain/project"
)

func (s *Server) setupRoutes() *http.ServeMux {
	connectrpcMux := http.NewServeMux()

	// Examples
	greeterHandler := greeter_domain.NewHandler(s.c.GetGreeterService())
	employeeHandler := employee_domain.NewHandler(s.c.GetEmployeeService())
	greeterPath, greeterConnectHandler := greeterv1connect.NewGreeterServiceHandler(greeterHandler)
	employeePath, employeeConnectHandler := altalunev1connect.NewEmployeeServiceHandler(employeeHandler)
	connectrpcMux.Handle(greeterPath, greeterConnectHandler)
	connectrpcMux.Handle(employeePath, employeeConnectHandler)

	// Domains
	projectHandler := project_domain.NewHandler(s.c.GetProjectService())
	projectPath, projectConnectHandler := altalunev1connect.NewProjectServiceHandler(projectHandler)
	connectrpcMux.Handle(projectPath, projectConnectHandler)

	apiKeyHandler := api_key_domain.NewHandler(s.c.GetApiKeyService())
	apiKeyPath, apiKeyConnectHandler := altalunev1connect.NewApiKeyServiceHandler(apiKeyHandler)
	connectrpcMux.Handle(apiKeyPath, apiKeyConnectHandler)

	// main server mux
	mux := http.NewServeMux()

	mux.Handle("/api/", http.StripPrefix("/api", connectrpcMux))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		health := map[string]any{
			"status": "ok",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(health); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	// serve frontend
	websiteFS, _ := fs.Sub(altalune.FrontendEmbeddedFiles, "frontend/.output/public")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
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

	return mux
}
