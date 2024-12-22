package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asolheiro/kiosk-api/internal/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phenpessoa/gutils/netutils/httputils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,
		os.Kill,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)

	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("Goodby...")
}

func run(ctx context.Context) error {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	
	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	logger = logger.Named("kiosk-api")
	defer func() {_ = logger.Sync() } ()

	pool, err := pgxpool.New(
		ctx,
		fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DB"),
		),
	)
	if err != nil {
		return err
	}; defer pool.Close()
	
	if err := pool.Ping(ctx); err != nil {
		return err
	}

	apiInstance := api.NewAPI(
		pool, 
		logger,
	)

	r := chi.NewMux()
	r.Use(
		middleware.RequestID,
		middleware.Recoverer,
		httputils.ChiLogger(logger),
	)

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		healthy := map[string]string{
			"status": "healthy",
		}
		jsonData, err := json.Marshal(healthy)
		if err != nil {
			http.Error(w, "error encoding json", http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})
	
	r.Route("/users", func(r chi.Router) {
		r.Post("/", apiInstance.PostUser)
		r.Get("/", apiInstance.ListUsers)
		r.Get("/{userId}", apiInstance.GetUser)
		r.Put("/{userId}", apiInstance.PutUser)
	})
	

	srv := http.Server{
		Addr: ":8080",
		Handler: r,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
	 }

	defer func () {
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(
			context.Background(), timeout,
		)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown server", zap.Error(err))
		}
	} ()

	errChan := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
		logger.Info("Starting server at port 8080...")
	} ()

	select {
	case <- ctx.Done():
		return nil
	case err := <- errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil 
}