package products

import (
	"context"

	"cloud.google.com/go/firestore"

	"basket-bliss-backend/internal/platform/database"
)

type Repository struct {
	db *firestore.Client
}

func NewRepository(db *firestore.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, merchantID string, p Product) (Product, error) {
	ref := database.Products(r.db, merchantID).NewDoc()
	if _, err := ref.Create(ctx, p); err != nil {
		return Product{}, err
	}
	p.ID = ref.ID
	return p, nil
}

func (r *Repository) List(ctx context.Context, merchantID string) ([]Product, error) {
	docs, err := database.Products(r.db, merchantID).OrderBy("createdAt", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]Product, 0, len(docs))
	for _, d := range docs {
		var p Product
		if err := d.DataTo(&p); err != nil {
			return nil, err
		}
		p.ID = d.Ref.ID
		out = append(out, p)
	}
	return out, nil
}
