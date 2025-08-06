package server

import (
	"bytes"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
)

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
