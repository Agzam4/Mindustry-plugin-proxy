package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type apiProxy struct {
	javaURL *url.URL
	rp      *httputil.ReverseProxy
}

func newAPIProxy(javaURLStr string) (*apiProxy, error) {
	ju, err := url.Parse(javaURLStr)
	if err != nil {
		return nil, err
	}
	p := &apiProxy{javaURL: ju}
	p.rp = &httputil.ReverseProxy{Director: p.direct}
	return p, nil
}

func (p *apiProxy) direct(req *http.Request) {
	req.URL.Scheme = p.javaURL.Scheme
	req.URL.Host = p.javaURL.Host
	req.Host = p.javaURL.Host
	req.URL.Path = strings.TrimPrefix(req.URL.Path, "/api")
	if !strings.HasPrefix(req.URL.Path, "/") {
		req.URL.Path = "/" + req.URL.Path
	}
}

func (p *apiProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.rp.ServeHTTP(w, r)
}
