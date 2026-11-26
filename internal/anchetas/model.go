package anchetas

import "time"

// Item es un producto dentro de la ancheta. ProductName es una copia del nombre
// para poder mostrar la ancheta sin leer cada producto.
type Item struct {
	ProductID   string `json:"productId" firestore:"productId"`
	ProductName string `json:"productName" firestore:"productName"`
	Quantity    int    `json:"quantity" firestore:"quantity"`
}

type Ancheta struct {
	ID          string    `json:"id" firestore:"-"`
	Name        string    `json:"name" firestore:"name"`
	Description string    `json:"description" firestore:"description"`
	Price       int64     `json:"price" firestore:"price"` // en pesos, sin decimales
	Items       []Item    `json:"items" firestore:"items"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`

	// Available es cuántas anchetas se pueden armar con el stock actual. Se calcula, no se guarda.
	Available int `json:"available" firestore:"-"`
}

type ItemInput struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type CreateInput struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       int64       `json:"price"`
	Items       []ItemInput `json:"items"`
}
