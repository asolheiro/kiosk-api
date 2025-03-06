package routers

import (
	"github.com/asolheiro/kiosk-api/internal-v2/controller"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func PrintRouter(server *fuego.Server, handler controller.APIResources) {
	print := fuego.Group(
		server,
		"/print",
		option.Summary("Print route"),
		option.Description("Default description for print route"),
		option.Tags("print"),
	)

	fuego.Post(print, "/", handler.PostPrint)
}