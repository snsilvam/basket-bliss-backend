package orders

import (
	"context"

	"basket-bliss-backend/internal/platform/apperr"
	"basket-bliss-backend/internal/platform/web"
)

const maxQuantity = 100

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create es el flujo 5: el customer compra una ancheta. No hay pasarela de pagos todavía.
func (s *Service) Create(ctx context.Context, merchantID, customerID string, in CreateInput) (Order, error) {
	if !web.ValidID(in.AnchetaID) {
		return Order{}, apperr.Invalid("anchetaId es obligatorio")
	}
	if in.Quantity < 1 || in.Quantity > maxQuantity {
		return Order{}, apperr.Invalid("la cantidad debe estar entre 1 y %d", maxQuantity)
	}
	return s.repo.Purchase(ctx, merchantID, customerID, in.AnchetaID, in.Quantity)
}

func (s *Service) List(ctx context.Context, merchantID, customerID string) ([]Order, error) {
	return s.repo.List(ctx, merchantID, customerID)
}
