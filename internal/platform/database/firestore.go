// Package database crea el cliente de Firestore y centraliza las rutas de las colecciones.
package database

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Connect abre el cliente de Firestore. Si projectID está vacío se detecta desde
// las credenciales. Si FIRESTORE_EMULATOR_HOST está definido, el cliente usa el emulador.
func Connect(ctx context.Context, projectID, databaseID string) (*firestore.Client, error) {
	if projectID == "" {
		projectID = firestore.DetectProjectID
	}
	return firestore.NewClientWithDatabase(ctx, projectID, databaseID)
}

// IsNotFound indica si un error de Firestore corresponde a un documento inexistente.
func IsNotFound(err error) bool {
	return status.Code(err) == codes.NotFound
}
