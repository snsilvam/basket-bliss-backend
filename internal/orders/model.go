package orders

import "time"

// Sin pasarela de pagos en v1: la orden queda creada y el stock ya descontado.
const StatusPendingPayment = "pending_payment"

type Order struct {
	ID          string    `json:"id" firestore:"-"`
	CustomerID  string    `json:"customerId" firestore:"customerId"`
	AnchetaID   string    `json:"anchetaId" firestore:"anchetaId"`
	AnchetaName string    `json:"anchetaName" firestore:"anchetaName"`
	Quantity    int       `json:"quantity" firestore:"quantity"`
	UnitPrice   int64     `json:"unitPrice" firestore:"unitPrice"`
	Total       int64     `json:"total" firestore:"total"`
	Status      string    `json:"status" firestore:"status"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
}

type CreateInput struct {
	AnchetaID string `json:"anchetaId"`
	Quantity  int    `json:"quantity"`
}
