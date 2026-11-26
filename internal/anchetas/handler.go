package anchetas

import (
	"net/http"

	"basket-bliss-backend/internal/platform/web"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/merchants/{merchantId}/anchetas", h.create)
	mux.HandleFunc("GET /api/merchants/{merchantId}/anchetas", h.list)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in CreateInput
	if err := web.Decode(w, r, &in); err != nil {
		web.Error(w, r, err)
		return
	}
	a, err := h.svc.Create(r.Context(), r.PathValue("merchantId"), in)
	if err != nil {
		web.Error(w, r, err)
		return
	}
	web.JSON(w, http.StatusCreated, a)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List(r.Context(), r.PathValue("merchantId"))
	if err != nil {
		web.Error(w, r, err)
		return
	}
	web.JSON(w, http.StatusOK, list)
}
