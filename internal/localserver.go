package awsserviceproxy

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

//go:embed frontend/dist/aws-web-proxy/*
var UI embed.FS

var uiFS fs.FS

func StartWebServer() {
	uiFS, _ = fs.Sub(UI, "frontend/dist/aws-web-proxy")
	go func() {
		rtr := mux.NewRouter()
		rtr.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(AWPConfig)
		}).Methods("GET")

		rtr.HandleFunc("/api/logs/{service:[a-z\\.\\-]+}", func(w http.ResponseWriter, r *http.Request) {
			params := mux.Vars(r)
			service := params["service"]
			pwd, _ := os.Getwd()
			logfile, err := ioutil.ReadFile(fmt.Sprintf("%s/logs/%s.log", pwd, service))
			w.Header().Set("Content-Type", "application/json")
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			logs := string(logfile)
			entries := remove(strings.Split(logs, "\n"), "")
			w.Write([]byte(fmt.Sprintf("[%s]", strings.Join(entries, ","))))
		}).Methods("GET")

		rtr.PathPrefix("/").HandlerFunc(handleStatic)
		log.Println("Access UI: http://localhost:2137")
		if err := http.ListenAndServe("localhost:2137", rtr); err != nil {
			log.Fatal("Failed. Try to execute `lsof -t -i tcp:2137 | xargs kill`.")
		}

	}()
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	path := filepath.Clean(r.URL.Path)
	if path == "/" { // Add other paths that you route on the UI side here
		path = "index.html"
	}
	path = strings.TrimPrefix(path, "/")

	file, err := uiFS.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("file", path, "not found:", err)
			http.NotFound(w, r)
			return
		}
		log.Println("file", path, "cannot be read:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", contentType)
	if strings.HasPrefix(path, "static/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	}
	stat, err := file.Stat()
	if err == nil && stat.Size() > 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}

	n, _ := io.Copy(w, file)
	log.Println("file", path, "copied", n, "bytes")
}
