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
	employee_domain "github.com/hrz8/altalune/pkg/employee"
	greeter_domain "github.com/hrz8/altalune/pkg/greeter"
	project_domain "github.com/hrz8/altalune/pkg/project"
)

func (s *Server) setupRoutes() *http.ServeMux {
	connectrpcMux := http.NewServeMux()
	greeterHandler := greeter_domain.NewHandler(s.c.GetGreeterService())
	employeeHandler := employee_domain.NewHandler(s.c.GetEmployeeService())
	projectHandler := project_domain.NewHandler(s.c.GetProjectService())

	greeterPath, greeterConnectHandler := greeterv1connect.NewGreeterServiceHandler(greeterHandler)
	employeePath, employeeConnectHandler := altalunev1connect.NewEmployeeServiceHandler(employeeHandler)
	projectPath, projectConnectHandler := altalunev1connect.NewProjectServiceHandler(projectHandler)

	connectrpcMux.Handle(greeterPath, greeterConnectHandler)
	connectrpcMux.Handle(employeePath, employeeConnectHandler)
	connectrpcMux.Handle(projectPath, projectConnectHandler)

	// main server mux
	mux := http.NewServeMux()

	mux.Handle("/api/", http.StripPrefix("/api", connectrpcMux))
	mux.HandleFunc("/healthz", s.healthCheckHandler)

	// serve frontend
	websiteFS, _ := fs.Sub(altalune.FrontendEmbeddedFiles, "frontend/.output/public")
	mux.HandleFunc("/", s.websiteHandler(websiteFS))

	return mux
}

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]any{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(health); err != nil {
		s.log.Error("failed to encode health check response", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) websiteHandler(websiteFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
