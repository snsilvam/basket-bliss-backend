package customers

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

func (r *Repository) Create(ctx context.Context, merchantID string, c Customer) (Customer, error) {
	ref := database.Customers(r.db, merchantID).NewDoc()
	if _, err := ref.Create(ctx, c); err != nil {
		return Customer{}, err
	}
	c.ID = ref.ID
	return c, nil
}

func (r *Repository) List(ctx context.Context, merchantID string) ([]Customer, error) {
	docs, err := database.Customers(r.db, merchantID).OrderBy("createdAt", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]Customer, 0, len(docs))
	for _, d := range docs {
		var c Customer
		if err := d.DataTo(&c); err != nil {
			return nil, err
		}
		c.ID = d.Ref.ID
		out = append(out, c)
	}
	return out, nil
}
