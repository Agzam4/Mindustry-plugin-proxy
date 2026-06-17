package main

import (
	"log"
	"net/http"
	"strings"
)

type corsWrap struct {
	next http.Handler
}

func (c *corsWrap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Vary", "Origin")
	}

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	c.next.ServeHTTP(w, r)
}

type wrapLog struct {
	next http.Handler
}

func (l *wrapLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.String())
	if strings.HasPrefix(r.URL.Path, "/api/") {
		log.Printf("  Headers: %s=%s, %s=%s",
			HeaderSID, r.Header.Get(HeaderSID),
			HeaderClientIP, r.Header.Get(HeaderClientIP))
	}
	l.next.ServeHTTP(w, r)
}
