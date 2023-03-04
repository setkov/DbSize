package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var staticFS embed.FS

type WebUI struct {
	config *Config
	server *http.Server
}

func NewWebUI(c *Config) *WebUI {
	return &WebUI{config: c}
}

func (w *WebUI) Start() {
	port := fmt.Sprintf(":%v", w.config.WebUI.Port)

	static, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Print(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(static)))
	mux.HandleFunc("/api/servers", w.Servers)
	mux.HandleFunc("/api/databases", w.Databases)

	w.server = &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("start WebUI at %v", port)
	go func() {
		w.server.ListenAndServe()
	}()
}

func (w *WebUI) Stop() {
	log.Print("stop WebUI")
	w.server.Shutdown(context.Background())
}
