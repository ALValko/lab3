package app

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

var (
	templatesDir = "app/templates"
	indexHTML    = path.Join(templatesDir, "index.html")
	baseHTML     = path.Join(templatesDir, "base.html")
	headerHTML   = path.Join(templatesDir, "header.html")
)

type handler struct {
	dir    string
	logger *log.Logger
}

func NewHandler(dir string, logger *log.Logger) *handler {
	return &handler{dir, logger}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := path.Join("files", r.URL.Path)
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		NotFoundError(w)
		return
	}

	if err != nil {
		h.logger.Printf("%s\n", err)
		InternalServerError(w)
		return
	}

	if stat.IsDir() {
		h.HandleDir(w, path, r.URL.Path)
	} else {
		h.HandleFile(w, path, stat.Name())
	}
}

type Result struct {
	Path    string
	Parent  string
	Entries []Entry
}

type Entry struct {
	Path       string
	Name       string
	IsDir      bool
	Size       int64
	ModifiedAt time.Time
	Mode       os.FileMode
}

func NewEntry(filePath string, info os.FileInfo) *Entry {
	return &Entry{
		Path:       path.Join(filePath, info.Name()),
		Name:       info.Name(),
		IsDir:      info.IsDir(),
		Size:       info.Size(),
		ModifiedAt: info.ModTime(),
		Mode:       info.Mode(),
	}
}

func (h *handler) HandleDir(w http.ResponseWriter, dirPath string, urlPath string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		h.logger.Printf("%s\n", err)
		InternalServerError(w)
		return
	}

	resultEntries := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err == nil {
			resultEntries = append(resultEntries, *NewEntry(urlPath, info))
		}
	}

	tmpl, err := template.ParseFiles(indexHTML, baseHTML, headerHTML)
	if err != nil {
		h.logger.Printf("%s\n", err)
		InternalServerError(w)
		return
	}

	result := &Result{
		Path:    path.Clean(urlPath),
		Parent:  path.Dir(path.Clean(urlPath)),
		Entries: resultEntries,
	}

	err = tmpl.Execute(w, result)
	if err != nil {
		h.logger.Printf("%s\n", err)
		InternalServerError(w)
	}
}

func (h *handler) HandleFile(w http.ResponseWriter, filePath string, fileName string) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		h.logger.Printf("%s\n", err)
		InternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Write(bytes)
}
