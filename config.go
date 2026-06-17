package main

import "os"

type Config struct {
	JavaAPIURL string
	ListenAddr string
	TLSMode    string
	TLSCert    string
	TLSKey     string
	Domain     string
	StaticDir  string
}

func LoadConfig() Config {
	return Config{
		JavaAPIURL: getEnv("JAVA_API_URL", "http://127.0.0.1:6568"),
		ListenAddr: getEnv("LISTEN_ADDR", ":8080"),
		TLSMode:    getEnv("TLS_MODE", "none"),
		TLSCert:    getEnv("TLS_CERT", ""),
		TLSKey:     getEnv("TLS_KEY", ""),
		Domain:     getEnv("DOMAIN", ""),
		StaticDir:  getEnv("STATIC_DIR", "./static"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
