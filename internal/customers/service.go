package customers

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

// Create es el flujo 3: registrar un customer del merchant.
func (s *Service) Create(ctx context.Context, merchantID string, in CreateInput) (Customer, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	if name == "" {
		return Customer{}, apperr.Invalid("el nombre es obligatorio")
	}
	if !strings.Contains(email, "@") {
		return Customer{}, apperr.Invalid("el email no es válido")
	}
	if _, err := s.merchants.Get(ctx, merchantID); err != nil {
		return Customer{}, err
	}
	return s.repo.Create(ctx, merchantID, Customer{Name: name, Email: email, CreatedAt: time.Now().UTC()})
}

func (s *Service) List(ctx context.Context, merchantID string) ([]Customer, error) {
	return s.repo.List(ctx, merchantID)
}
