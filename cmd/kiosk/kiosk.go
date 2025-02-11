package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asolheiro/kiosk-api/internal/api"
	"github.com/asolheiro/kiosk-api/internal/utils"
	_ "github.com/asolheiro/kiosk-api/docs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/phenpessoa/gutils/netutils/httputils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/swaggo/http-swagger"
	_ "modernc.org/sqlite"
)

func main() {
	fmt.Println("Starting kiosk-api...")
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
	fmt.Println("Goodbye...")
}

var ddl string

func run(ctx context.Context) error {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	logger = logger.Named("kiosk-api")
	defer func() { _ = logger.Sync() }()

	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		return fmt.Errorf("SQLITE_DB_PATH environment variable is not set")
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	apiInstance := api.NewSQliteAPI(
		db,
		logger,
	)

	r := chi.NewMux()
	r.Use(
		middleware.RequestID,
		middleware.Recoverer,
		httputils.ChiLogger(logger),
	)


	r.Route("/", func(r chi.Router) {
		r.Get("/healthcheck", utils.HealthCheck)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	utils.UsersRouter(r, apiInstance)
	utils.EventsRouter(r, apiInstance)
	utils.GuestsRouter(r, apiInstance)
	utils.CheckinsRouter(r, apiInstance)
	utils.ConfigRouter(r, apiInstance)
	utils.PrintRouter(r, apiInstance)


	srv := http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	defer func() {
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(
			context.Background(), timeout,
		)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Failed to shutdown server", zap.Error(err))
		}
	}()

	errChan := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
		logger.Info("Starting server at port 8080...")
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}
