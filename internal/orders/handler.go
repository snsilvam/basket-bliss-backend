package orders

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
	mux.HandleFunc("POST /api/merchants/{merchantId}/customers/{customerId}/orders", h.create)
	mux.HandleFunc("GET /api/merchants/{merchantId}/customers/{customerId}/orders", h.list)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in CreateInput
	if err := web.Decode(w, r, &in); err != nil {
		web.Error(w, r, err)
		return
	}
	o, err := h.svc.Create(r.Context(), r.PathValue("merchantId"), r.PathValue("customerId"), in)
	if err != nil {
		web.Error(w, r, err)
		return
	}
	web.JSON(w, http.StatusCreated, o)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List(r.Context(), r.PathValue("merchantId"), r.PathValue("customerId"))
	if err != nil {
		web.Error(w, r, err)
		return
	}
	web.JSON(w, http.StatusOK, list)
}
