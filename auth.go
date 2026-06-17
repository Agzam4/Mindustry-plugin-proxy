package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type authCallback struct {
	javaBase string
}

type authResp struct {
	ID    string `json:"id"`
	UUID  string `json:"uuid"`
	Error string `json:"error"`
}

func (h *authCallback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token parameter", http.StatusBadRequest)
		return
	}

	body := fmt.Sprintf(`{"token":"%s"}`, token)
	req, err := http.NewRequest(http.MethodPost,
		h.javaBase+"/auth/create-session",
		strings.NewReader(body))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(HeaderClientIP, clientIP(r))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "auth backend unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var ar authResp
	json.Unmarshal(raw, &ar)

	if ar.Error != "" || ar.ID == "" {
		http.Error(w, "authentication failed: "+ar.Error, http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    ar.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.TLSMode != "none",
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
