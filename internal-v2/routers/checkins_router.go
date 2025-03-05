package routers

import (
	"github.com/asolheiro/kiosk-api/internal-v2/controller"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func CheckInsRouter(server *fuego.Server, handler controller.APIResources) {
	checkin := fuego.Group(
		server,
		"/checkin",
		option.Summary("CheckIn routes"),
		option.Description("Default description for all CheckIn routes"),
		option.Tags("checkin"),
	)

	fuego.Post(checkin, "/", handler.CreateCheckIn,
		option.Summary("Create a CheckIn"),
		option.Description("Create a CheckIn with the given data"),
	)
	
	fuego.Get(checkin, "/", handler.ListCheckIns,
		option.Summary("Get all CheckIn"),
		option.Description("Returns a list of all CheckIns"),
	)
	
	fuego.Get(checkin, "/{checkInId}", handler.FindCheckInById,
		option.Summary("Get a single CheckIn"),
		option.Description("Returns a CheckIn for the given Idlist of all events"),
	)
}