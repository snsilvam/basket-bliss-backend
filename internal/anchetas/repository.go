package anchetas

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

func (r *Repository) Create(ctx context.Context, merchantID string, a Ancheta) (Ancheta, error) {
	ref := database.Anchetas(r.db, merchantID).NewDoc()
	if _, err := ref.Create(ctx, a); err != nil {
		return Ancheta{}, err
	}
	a.ID = ref.ID
	return a, nil
}

func (r *Repository) List(ctx context.Context, merchantID string) ([]Ancheta, error) {
	docs, err := database.Anchetas(r.db, merchantID).OrderBy("createdAt", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]Ancheta, 0, len(docs))
	for _, d := range docs {
		var a Ancheta
		if err := d.DataTo(&a); err != nil {
			return nil, err
		}
		a.ID = d.Ref.ID
		out = append(out, a)
	}
	return out, nil
}
