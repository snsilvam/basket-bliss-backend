// Package config lee la configuración del servicio desde variables de entorno.
package config

import "os"

type Config struct {
	Host       string // vacío = todas las interfaces (lo que espera Cloud Run)
	Port       string
	ProjectID  string // vacío = se detecta desde las credenciales (ADC)
	DatabaseID string
}

func Load() Config {
	return Config{
		Host:       os.Getenv("HOST"),
		Port:       getenv("PORT", "8080"),
		ProjectID:  os.Getenv("GCP_PROJECT_ID"),
		DatabaseID: getenv("FIRESTORE_DATABASE_ID", "(default)"),
	}
}

func (c Config) Addr() string { return c.Host + ":" + c.Port }

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
