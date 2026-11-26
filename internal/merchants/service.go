package merchants

import (
	"context"
	"strings"
	"time"

	"basket-bliss-backend/internal/platform/apperr"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create es el flujo 1: onboarding del merchant.
func (s *Service) Create(ctx context.Context, in CreateInput) (Merchant, error) {
	name := strings.TrimSpace(in.Name)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	if name == "" {
		return Merchant{}, apperr.Invalid("el nombre es obligatorio")
	}
	if !strings.Contains(email, "@") {
		return Merchant{}, apperr.Invalid("el email no es válido")
	}
	return s.repo.Create(ctx, Merchant{Name: name, Email: email, CreatedAt: time.Now().UTC()})
}

func (s *Service) Get(ctx context.Context, id string) (Merchant, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Merchant, error) {
	return s.repo.List(ctx)
}
