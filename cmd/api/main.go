// Punto de entrada: carga la configuración, conecta Firestore, arma las dependencias y arranca el HTTP.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"basket-bliss-backend/internal/anchetas"
	"basket-bliss-backend/internal/customers"
	"basket-bliss-backend/internal/merchants"
	"basket-bliss-backend/internal/orders"
	"basket-bliss-backend/internal/platform/config"
	"basket-bliss-backend/internal/platform/database"
	"basket-bliss-backend/internal/platform/web"
	"basket-bliss-backend/internal/products"
)

func main() {
	if err := run(); err != nil {
		slog.Error("el servidor terminó con error", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	db, err := database.Connect(ctx, cfg.ProjectID, cfg.DatabaseID)
	if err != nil {
		return fmt.Errorf("conectando a Firestore: %w", err)
	}
	defer db.Close()

	merchantSvc := merchants.NewService(merchants.NewRepository(db))
	productSvc := products.NewService(products.NewRepository(db), merchantSvc)
	anchetaSvc := anchetas.NewService(anchetas.NewRepository(db), productSvc)
	customerSvc := customers.NewService(customers.NewRepository(db), merchantSvc)
	orderSvc := orders.NewService(orders.NewRepository(db))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		web.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	merchants.NewHandler(merchantSvc).Register(mux)
	products.NewHandler(productSvc).Register(mux)
	anchetas.NewHandler(anchetaSvc).Register(mux)
	customers.NewHandler(customerSvc).Register(mux)
	orders.NewHandler(orderSvc).Register(mux)

	srv := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           web.LogRequests(web.RejectEncodedSlashes(mux)),
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("API escuchando",
			"addr", cfg.Addr(),
			"project", cfg.ProjectID,
			"database", cfg.DatabaseID,
			"emulator", os.Getenv("FIRESTORE_EMULATOR_HOST"))
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
	}

	slog.Info("apagando el servidor")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
