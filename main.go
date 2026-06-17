package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

var cfg = LoadConfig()

func main() {
	ap, err := newAPIProxy(cfg.JavaAPIURL)
	if err != nil {
		log.Fatalf("proxy: %v", err)
	}

	auth := &authCallback{javaBase: cfg.JavaAPIURL}
	spa := &spaFS{root: http.Dir(cfg.StaticDir), indexPath: cfg.StaticDir}

	mux := http.NewServeMux()
	mux.Handle("/api/", &sessionWare{next: ap})
	mux.Handle("/", spa)
	mux.Handle("/auth/callback", auth)

	var h http.Handler = mux
	if cfg.TLSMode == "none" {
		h = &corsWrap{next: h}
	}
	h = &wrapLog{next: h}

	switch cfg.TLSMode {
	case "none":
		log.Printf("HTTP on %s", cfg.ListenAddr)
		log.Fatal(http.ListenAndServe(cfg.ListenAddr, h))

	case "cert":
		log.Printf("HTTPS on %s (cert)", cfg.ListenAddr)
		log.Fatal(http.ListenAndServeTLS(cfg.ListenAddr, cfg.TLSCert, cfg.TLSKey, h))

	case "auto":
		log.Printf("HTTPS on %s (auto/Let's Encrypt)", cfg.ListenAddr)
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(cfg.Domain),
			Cache:      autocert.DirCache("certs"),
		}
		srv := &http.Server{
			Addr:    cfg.ListenAddr,
			Handler: h,
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
			},
		}
		log.Fatal(srv.ListenAndServeTLS("", ""))

	default:
		log.Fatalf("unknown TLS_MODE: %s (expected none|cert|auto)", cfg.TLSMode)
	}
}

type spaFS struct {
	root      http.FileSystem
	indexPath string
}

func (s *spaFS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	path := filepath.Join(s.indexPath, r.URL.Path)
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		r.URL.Path = "/"
		http.FileServer(s.root).ServeHTTP(w, r)
		return
	}
	http.FileServer(s.root).ServeHTTP(w, r)
}
