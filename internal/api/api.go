package api

import (
	"github.com/asolheiro/kiosk-api/internal/pgstore"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type API struct {
	repo *pgstore.Queries
	logger *zap.Logger
	validator *validator.Validate
	pool *pgxpool.Pool
}

func NewAPI(pool *pgxpool.Pool, logger *zap.Logger) API {
	validator := validator.New(validator.WithRequiredStructEnabled())

	return API{
		pgstore.New(pool),
		logger,
		validator,
		pool,
	}
}