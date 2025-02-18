package api

import (
	"database/sql"

	"github.com/asolheiro/kiosk-api/internal/sqlitestore"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type API struct {
	repo      *sqlitestore.Queries
	logger    *zap.Logger
	validator *validator.Validate
	pool      *sql.DB
}

func NewSQliteAPI(db *sql.DB, logger *zap.Logger) API {
	validator := validator.New(validator.WithRequiredStructEnabled())

	return API{
		sqlitestore.New(db),
		logger,
		validator,
		db,
	}
}
