package products

import "time"

// Product es un producto físico del merchant. Su stock es el inventario real y
// se descuenta al vender una ancheta que lo contiene.
type Product struct {
	ID        string    `json:"id" firestore:"-"`
	Name      string    `json:"name" firestore:"name"`
	Stock     int       `json:"stock" firestore:"stock"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
}

type CreateInput struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}
