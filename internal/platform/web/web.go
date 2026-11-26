// Package web tiene los helpers HTTP compartidos: JSON, errores y log de peticiones.
package web

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"basket-bliss-backend/internal/platform/apperr"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("no se pudo escribir la respuesta", "err", err)
	}
}

func Decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperr.Invalid("el cuerpo de la petición no es un JSON válido")
	}
	return nil
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	var ae *apperr.Error
	if errors.As(err, &ae) {
		status := http.StatusBadRequest
		switch ae.Kind {
		case apperr.KindNotFound:
			status = http.StatusNotFound
		case apperr.KindConflict:
			status = http.StatusConflict
		}
		JSON(w, status, map[string]string{"error": ae.Message})
		return
	}
	slog.Error("error interno", "method", r.Method, "path", r.URL.Path, "err", err)
	JSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno del servidor"})
}

// ValidID rechaza IDs vacíos o con "/" (en Firestore la barra separa segmentos de la ruta).
func ValidID(id string) bool {
	return id != "" && !strings.Contains(id, "/")
}

// RejectEncodedSlashes bloquea rutas con "%2F": el ID ya decodificado traería una "/"
// y apuntaría a otra ruta de Firestore.
func RejectEncodedSlashes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(strings.ToLower(r.URL.RawPath), "%2f") {
			JSON(w, http.StatusBadRequest, map[string]string{"error": "ruta inválida"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		slog.Info("http", "method", r.Method, "path", r.URL.Path, "status", rec.status, "duration", time.Since(start).Round(time.Millisecond))
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
