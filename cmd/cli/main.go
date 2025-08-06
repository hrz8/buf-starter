package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"github.com/hrz8/starter"
	greeterv1 "github.com/hrz8/starter/gen/greeter/v1"
	"github.com/hrz8/starter/gen/greeter/v1/greeterv1connect"
)

type Service struct {
	greeterv1.UnimplementedGreeterServiceServer

	validator protovalidate.Validator
}

func NewService(v protovalidate.Validator) *Service {
	return &Service{
		validator: v,
	}
}

var allowedNameMap = map[string]bool{
	"Alina": true, "Bryce": true, "Carmen": true, "Darius": true, "Elena": true,
	"Felix": true, "Gianna": true, "Hassan": true, "Irene": true, "Jasper": true,
	"Kiana": true, "Luther": true, "Maya": true, "Nolan": true, "Orlando": true,
	"Priya": true, "Quincy": true, "Rafael": true, "Sienna": true, "Tobias": true,
	"Umair": true, "Vera": true, "Wesley": true, "Xavier": true, "Yasmin": true,
	"Zane": true, "Adriana": true, "Bennett": true, "Clarissa": true, "Devonte": true,
	"Estella": true, "Finnegan": true, "Gracelyn": true, "Harvey": true, "Isidora": true,
	"Jovani": true, "Katarina": true, "Leonidas": true, "Mirella": true, "Nikolas": true,
	"Octavia": true, "Percival": true, "Quintessa": true, "Romero": true, "Salvador": true,
	"Theodora": true, "Ulrich": true, "Valeria": true, "Winslow": true, "Xiomara": true,
	"Yuridia": true, "Zephyrus": true, "Aurelius": true, "Bellatrix": true, "Caspian": true,
	"Demetrius": true, "Evangeline": true, "Florentino": true, "Galadriel": true, "Hermione": true,
	"Ignatius": true, "Julianna": true, "Kristoffer": true, "Lysandra": true, "Maximiliano": true,
	"Nefertari": true, "Olivander": true, "Philomena": true, "Quetzalcoatl": true, "Rhiannon": true,
	"Sebastiana": true, "Thessalonia": true, "Ulyssiana": true, "Vladimir": true, "Wilhelmina": true,
	"Xenophilius": true, "Yggdrasila": true, "Zaphkiel": true, "Alejandrina": true, "Balthazar": true,
	"Christabelle": true, "Domenico": true, "Euphrosyne": true, "Featherstone": true, "Gwendolyn": true,
	"Hyacinthus": true, "Isambard": true, "Jacqueline": true, "Kallistrate": true, "Leontius": true,
	"Marcellinus": true, "Nicomachus": true, "Ozymandias": true, "Petronella": true, "Quintilius": true,
	"Rosencrantz": true, "Seraphimiel": true, "Timotheus": true, "Ultraviolet": true, "Valentinian": true,
}

func (s *Service) SayHello(ctx context.Context, req *greeterv1.SayHelloRequest) (*greeterv1.SayHelloResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, starter.NewInvalidPayloadError(err.Error())
	}

	if !allowedNameMap[req.Name] {
		return nil, starter.NewGreetingUnrecognize(req.Name)
	}

	response := &greeterv1.SayHelloResponse{
		Message: "Hello, " + req.Name,
	}
	return response, nil
}

type Handler struct {
	svc greeterv1.GreeterServiceServer
}

func NewHandler(svc greeterv1.GreeterServiceServer) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) SayHello(
	ctx context.Context,
	req *connect.Request[greeterv1.SayHelloRequest],
) (*connect.Response[greeterv1.SayHelloResponse], error) {
	response, err := h.svc.SayHello(ctx, req.Msg)
	if err != nil {
		return nil, starter.ToConnectError(err)
	}
	return connect.NewResponse(response), nil
}

type Server struct {
	httpserver *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Bootstrap(
	greeterHandler greeterv1connect.GreeterServiceHandler,
) {
	connectmux := http.NewServeMux()

	greeterPath, greeterConnect := greeterv1connect.NewGreeterServiceHandler(greeterHandler)
	connectmux.Handle(greeterPath, greeterConnect)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", connectmux))

	websiteFS, _ := fs.Sub(starter.FrontendEmbeddedFiles, "frontend/.output/public")
	mux.HandleFunc("/", s.websiteHandler(websiteFS))

	s.httpserver = &http.Server{
		Addr:    fmt.Sprintf(":%d", 3000),
		Handler: mux,
	}
}

func (s *Server) Start() {
	if s.httpserver == nil {
		panic("Server not initialized. Call Bootstrap first.")
	}

	fmt.Println("Starting server on", s.httpserver.Addr)
	if err := s.httpserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func main() {
	v, _ := protovalidate.New()
	svc := NewService(v)
	handler := NewHandler(svc)
	server := NewServer()
	server.Bootstrap(handler)

	server.Start()
}

// websiteHandler serving frontend website application page
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

func exists(fsys fs.FS, name string) bool {
	_, err := fs.Stat(fsys, name)
	return err == nil
}

func isDir(fsys fs.FS, name string) bool {
	fi, err := fs.Stat(fsys, name)
	return err == nil && fi.IsDir()
}

func serveFileOr404(w http.ResponseWriter, r *http.Request, fsys fs.FS, name string) {
	f, err := fsys.Open(name)
	if err != nil {
		serve404Page(w, r, fsys)
		return
	}
	defer f.Close()

	fi, err := fs.Stat(fsys, name)
	if err != nil {
		serve404Page(w, r, fsys)
		return
	}

	ext := filepath.Ext(name)
	if ctype := mime.TypeByExtension(ext); ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	http.ServeContent(w, r, name, fi.ModTime(), bytes.NewReader(mustReadAll(f)))
}

func serve404Page(w http.ResponseWriter, r *http.Request, fsys fs.FS) {
	const notFound = "404.html"

	f, err := fsys.Open(notFound)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	data, _ := io.ReadAll(f)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write(data)
}

func mustReadAll(f fs.File) []byte {
	data, _ := io.ReadAll(f)
	return data
}
