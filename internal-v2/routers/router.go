package routers

import (
	"database/sql"

	"github.com/asolheiro/kiosk-api/internal-v2/controller"
	"github.com/go-fuego/fuego"
	"go.uber.org/zap"
)

func NewRouter(server *fuego.Server, db *sql.DB, logger *zap.Logger) {
	resources := controller.NewResource(db, logger)

	UsersRouter(server, *resources)
	EventsRouter(server, *resources)
	GuestRouter(server, *resources)
	CheckInsRouter(server, *resources)
}
