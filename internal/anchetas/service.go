package anchetas

import (
	"context"
	"strings"
	"time"

	"basket-bliss-backend/internal/platform/apperr"
	"basket-bliss-backend/internal/products"
)

type productLister interface {
	List(ctx context.Context, merchantID string) ([]products.Product, error)
}

type Service struct {
	repo     *Repository
	products productLister
}

func NewService(repo *Repository, products productLister) *Service {
	return &Service{repo: repo, products: products}
}

// Create es el flujo 2: registrar una ancheta con productos que ya existen.
func (s *Service) Create(ctx context.Context, merchantID string, in CreateInput) (Ancheta, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return Ancheta{}, apperr.Invalid("el nombre de la ancheta es obligatorio")
	}
	if in.Price <= 0 {
		return Ancheta{}, apperr.Invalid("el precio debe ser mayor que cero")
	}
	if len(in.Items) == 0 {
		return Ancheta{}, apperr.Invalid("la ancheta debe tener al menos un producto")
	}

	stock, names, err := s.stockByProduct(ctx, merchantID)
	if err != nil {
		return Ancheta{}, err
	}

	items := make([]Item, 0, len(in.Items))
	seen := make(map[string]bool, len(in.Items))
	for _, it := range in.Items {
		productName, ok := names[it.ProductID]
		if !ok {
			return Ancheta{}, apperr.Invalid("el producto %q no existe", it.ProductID)
		}
		if seen[it.ProductID] {
			return Ancheta{}, apperr.Invalid("el producto %q está repetido", productName)
		}
		if it.Quantity < 1 {
			return Ancheta{}, apperr.Invalid("la cantidad de %q debe ser al menos 1", productName)
		}
		seen[it.ProductID] = true
		items = append(items, Item{ProductID: it.ProductID, ProductName: productName, Quantity: it.Quantity})
	}

	a, err := s.repo.Create(ctx, merchantID, Ancheta{
		Name:        name,
		Description: strings.TrimSpace(in.Description),
		Price:       in.Price,
		Items:       items,
		CreatedAt:   time.Now().UTC(),
	})
	if err != nil {
		return Ancheta{}, err
	}
	a.Available = available(a.Items, stock)
	return a, nil
}

// List es el flujo 4: ver las anchetas en venta con su disponibilidad actual.
func (s *Service) List(ctx context.Context, merchantID string) ([]Ancheta, error) {
	list, err := s.repo.List(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	stock, _, err := s.stockByProduct(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i].Available = available(list[i].Items, stock)
	}
	return list, nil
}

func (s *Service) stockByProduct(ctx context.Context, merchantID string) (stock map[string]int, names map[string]string, err error) {
	prods, err := s.products.List(ctx, merchantID)
	if err != nil {
		return nil, nil, err
	}
	stock = make(map[string]int, len(prods))
	names = make(map[string]string, len(prods))
	for _, p := range prods {
		stock[p.ID] = p.Stock
		names[p.ID] = p.Name
	}
	return stock, names, nil
}

// available devuelve cuántas anchetas completas alcanzan con el stock: el mínimo de stock/cantidad.
func available(items []Item, stock map[string]int) int {
	if len(items) == 0 {
		return 0
	}
	n := -1
	for _, it := range items {
		if it.Quantity <= 0 {
			continue
		}
		c := stock[it.ProductID] / it.Quantity
		if n == -1 || c < n {
			n = c
		}
	}
	if n < 0 {
		return 0
	}
	return n
}
