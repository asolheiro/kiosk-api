package controller

import (
	"database/sql"

	"github.com/asolheiro/kiosk-api/internal-v2/sqlitestore"
	"go.uber.org/zap"
)

type APIResources struct {
	pool *sql.DB
	repo *sqlitestore.Queries
	logger *zap.Logger
}

func NewResource(db *sql.DB, logger *zap.Logger) *APIResources {
	return &APIResources{
		db,
		sqlitestore.New(db),
		logger,
	}
}

