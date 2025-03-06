package routers

import (
	"github.com/asolheiro/kiosk-api/internal-v2/controller"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func GuestRouter(server *fuego.Server, handler controller.APIResources) {
	guest := fuego.Group(
		server, 
		"/guests",
		option.Summary("Guests routes"),
		option.Description("Default description for all Guest routes"),
		option.Tags("guests"),
	)
	fuego.Post(guest, "/", handler.CreateGuest,
		option.Summary("Create a new guest"),
		option.Description("Crete a new guest with givern parameters"),
	)

	fuego.Get(guest, "/", handler.ListGuests, 
		option.Summary("Get all events"),
		option.Description("Returns a list of all events"),
	)

	fuego.Get(guest, "/{guestId}", handler.FindGuestById,
		option.Summary("Get a specific guest"),
		option.Description("Returns the guest with the given Id"),
	)

	fuego.Get(guest, "/document/{guestDocument}", handler.FindGuestByDocument,
		option.Summary("Get a specific guest"),
		option.Description("Returns the guest with the given document number"),
	)

	fuego.Put(guest, "/{guestId}", handler.UpdateGuest,
		option.Summary("Update a guest"),
		option.Description("Update a guest with the given parameters"),
	)
	fuego.Delete(guest, "/{guestId}", handler.DeleteGuest,
		option.Summary("Delete a event"),
		option.Description("Soft delete a event"),)
}
