package database

import "cloud.google.com/go/firestore"

// Estructura en Firestore (sigue el modelo de datos de la arquitectura v1):
//
//	merchants/{merchantId}
//	merchants/{merchantId}/products/{productId}
//	merchants/{merchantId}/anchetas/{anchetaId}
//	merchants/{merchantId}/customers/{customerId}
//	merchants/{merchantId}/customers/{customerId}/orders/{orderId}

func Merchants(db *firestore.Client) *firestore.CollectionRef {
	return db.Collection("merchants")
}

func Products(db *firestore.Client, merchantID string) *firestore.CollectionRef {
	return Merchants(db).Doc(merchantID).Collection("products")
}

func Anchetas(db *firestore.Client, merchantID string) *firestore.CollectionRef {
	return Merchants(db).Doc(merchantID).Collection("anchetas")
}

func Customers(db *firestore.Client, merchantID string) *firestore.CollectionRef {
	return Merchants(db).Doc(merchantID).Collection("customers")
}

func Orders(db *firestore.Client, merchantID, customerID string) *firestore.CollectionRef {
	return Customers(db, merchantID).Doc(customerID).Collection("orders")
}
