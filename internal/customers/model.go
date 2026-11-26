package customers

import "time"

type Customer struct {
	ID        string    `json:"id" firestore:"-"`
	Name      string    `json:"name" firestore:"name"`
	Email     string    `json:"email" firestore:"email"`
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
}

type CreateInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
