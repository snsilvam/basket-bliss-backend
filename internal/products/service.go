package products

import (
	"context"
	"strings"
	"time"

	"basket-bliss-backend/internal/merchants"
	"basket-bliss-backend/internal/platform/apperr"
)

type merchantGetter interface {
	Get(ctx context.Context, id string) (merchants.Merchant, error)
}

type Service struct {
	repo      *Repository
	merchants merchantGetter
}

func NewService(repo *Repository, merchants merchantGetter) *Service {
	return &Service{repo: repo, merchants: merchants}
}

func (s *Service) Create(ctx context.Context, merchantID string, in CreateInput) (Product, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return Product{}, apperr.Invalid("el nombre del producto es obligatorio")
	}
	if in.Stock < 0 {
		return Product{}, apperr.Invalid("el stock no puede ser negativo")
	}
	if _, err := s.merchants.Get(ctx, merchantID); err != nil {
		return Product{}, err
	}
	return s.repo.Create(ctx, merchantID, Product{Name: name, Stock: in.Stock, CreatedAt: time.Now().UTC()})
}

func (s *Service) List(ctx context.Context, merchantID string) ([]Product, error) {
	return s.repo.List(ctx, merchantID)
}
