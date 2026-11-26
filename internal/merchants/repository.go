package merchants

import (
	"context"

	"cloud.google.com/go/firestore"

	"basket-bliss-backend/internal/platform/apperr"
	"basket-bliss-backend/internal/platform/database"
)

type Repository struct {
	db *firestore.Client
}

func NewRepository(db *firestore.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, m Merchant) (Merchant, error) {
	ref := database.Merchants(r.db).NewDoc()
	if _, err := ref.Create(ctx, m); err != nil {
		return Merchant{}, err
	}
	m.ID = ref.ID
	return m, nil
}

func (r *Repository) Get(ctx context.Context, id string) (Merchant, error) {
	snap, err := database.Merchants(r.db).Doc(id).Get(ctx)
	if database.IsNotFound(err) {
		return Merchant{}, apperr.NotFound("el merchant %q no existe", id)
	}
	if err != nil {
		return Merchant{}, err
	}
	var m Merchant
	if err := snap.DataTo(&m); err != nil {
		return Merchant{}, err
	}
	m.ID = snap.Ref.ID
	return m, nil
}

func (r *Repository) List(ctx context.Context) ([]Merchant, error) {
	docs, err := database.Merchants(r.db).OrderBy("createdAt", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]Merchant, 0, len(docs))
	for _, d := range docs {
		var m Merchant
		if err := d.DataTo(&m); err != nil {
			return nil, err
		}
		m.ID = d.Ref.ID
		out = append(out, m)
	}
	return out, nil
}
