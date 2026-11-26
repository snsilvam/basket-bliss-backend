package orders

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"

	"basket-bliss-backend/internal/anchetas"
	"basket-bliss-backend/internal/platform/apperr"
	"basket-bliss-backend/internal/platform/database"
	"basket-bliss-backend/internal/products"
)

type Repository struct {
	db *firestore.Client
}

func NewRepository(db *firestore.Client) *Repository {
	return &Repository{db: db}
}

// Purchase hace la compra en una sola transacción: comprueba el stock de cada
// producto de la ancheta, lo descuenta y crea la orden. Si dos compras compiten
// por el mismo stock, Firestore reintenta la transacción y no hay sobreventa.
func (r *Repository) Purchase(ctx context.Context, merchantID, customerID, anchetaID string, qty int) (Order, error) {
	var order Order
	err := r.db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// En una transacción de Firestore todas las lecturas van antes que las escrituras.
		if _, err := tx.Get(database.Customers(r.db, merchantID).Doc(customerID)); err != nil {
			if database.IsNotFound(err) {
				return apperr.NotFound("el customer %q no existe", customerID)
			}
			return err
		}

		snap, err := tx.Get(database.Anchetas(r.db, merchantID).Doc(anchetaID))
		if err != nil {
			if database.IsNotFound(err) {
				return apperr.NotFound("la ancheta %q no existe", anchetaID)
			}
			return err
		}
		var a anchetas.Ancheta
		if err := snap.DataTo(&a); err != nil {
			return err
		}

		refs := make([]*firestore.DocumentRef, len(a.Items))
		for i, it := range a.Items {
			refs[i] = database.Products(r.db, merchantID).Doc(it.ProductID)
		}
		snaps, err := tx.GetAll(refs)
		if err != nil {
			return err
		}

		newStock := make([]int, len(a.Items))
		for i, it := range a.Items {
			if !snaps[i].Exists() {
				return apperr.Conflict("el producto %q ya no existe", it.ProductName)
			}
			var p products.Product
			if err := snaps[i].DataTo(&p); err != nil {
				return err
			}
			need := it.Quantity * qty
			if p.Stock < need {
				return apperr.Conflict("no hay stock suficiente de %q: hay %d y se necesitan %d", p.Name, p.Stock, need)
			}
			newStock[i] = p.Stock - need
		}

		for i, ref := range refs {
			if err := tx.Update(ref, []firestore.Update{{Path: "stock", Value: newStock[i]}}); err != nil {
				return err
			}
		}

		ref := database.Orders(r.db, merchantID, customerID).NewDoc()
		order = Order{
			CustomerID:  customerID,
			AnchetaID:   anchetaID,
			AnchetaName: a.Name,
			Quantity:    qty,
			UnitPrice:   a.Price,
			Total:       a.Price * int64(qty),
			Status:      StatusPendingPayment,
			CreatedAt:   time.Now().UTC(),
		}
		if err := tx.Create(ref, order); err != nil {
			return err
		}
		order.ID = ref.ID
		return nil
	})
	if err != nil {
		return Order{}, err
	}
	return order, nil
}

func (r *Repository) List(ctx context.Context, merchantID, customerID string) ([]Order, error) {
	docs, err := database.Orders(r.db, merchantID, customerID).OrderBy("createdAt", firestore.Desc).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	out := make([]Order, 0, len(docs))
	for _, d := range docs {
		var o Order
		if err := d.DataTo(&o); err != nil {
			return nil, err
		}
		o.ID = d.Ref.ID
		out = append(out, o)
	}
	return out, nil
}
