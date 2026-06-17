package main

import (
	"net"
	"net/http"
	"strings"
)

const (
	SessionCookie  = "session"
	HeaderSID      = "Agzam4-Session-Id"
	HeaderClientIP = "Agzam4-Client-Ip"
)

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if h := firstComma(xff); h != "" {
			return h
		}
	}
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func firstComma(s string) string {
	for i := range s {
		if s[i] == ',' {
			return strings.TrimSpace(s[:i])
		}
	}
	return strings.TrimSpace(s)
}

type sessionWare struct {
	next http.Handler
}

func (m *sessionWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sid := ""
	if c, err := r.Cookie(SessionCookie); err == nil && c != nil {
		sid = c.Value
	}
	r.Header.Set(HeaderSID, sid)
	r.Header.Set(HeaderClientIP, clientIP(r))
	m.next.ServeHTTP(w, r)
}
