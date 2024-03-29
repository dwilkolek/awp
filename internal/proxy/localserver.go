package proxy

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

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tfmcdigital/aws-web-proxy/internal/domain"
	"github.com/tfmcdigital/aws-web-proxy/internal/proxy/logger"
	"github.com/tfmcdigital/aws-web-proxy/internal/proxy/websocket"
	"github.com/tfmcdigital/aws-web-proxy/internal/utils"
)

//go:embed frontend/dist/aws-web-proxy/*
var UI embed.FS

var uiFS fs.FS

func localWebServer() {

	uiFS, _ = fs.Sub(UI, "frontend/dist/aws-web-proxy")

	rtr := mux.NewRouter()

	rtr.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(domain.GetConfig())
	}).Methods("GET")

	rtr.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(w, r)
	})

	rtr.HandleFunc("/api/logs/{service:[a-z\\.\\-]+}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		service := params["service"]
		data, err := logs(service)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(data))
	}).Methods("GET")

	rtr.PathPrefix("/").HandlerFunc(handleStatic)
	log.Printf("Access UI: http://awp and http://localhost:%d", WEB_SERVER_PORT)

	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", WEB_SERVER_PORT), handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), handlers.AllowCredentials(),
	)(rtr)); err != nil {
		log.Fatal("Failed. Try to execute `lsof -t -i tcp:2137 | xargs kill`.")
	}

}

func logs(service string) (string, error) {
	logfile, err := ioutil.ReadFile(logger.GetLogFileLocation(service))
	if err != nil {
		return "", nil
	}
	logs := string(logfile)
	entries := utils.Remove(strings.Split(logs, "\n"), "")
	return fmt.Sprintf("[%s]", strings.Join(entries, ",")), err
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

	io.Copy(w, file)
}
