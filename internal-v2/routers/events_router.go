package routers

import (
	"github.com/asolheiro/kiosk-api/internal-v2/controller"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func EventsRouter(server *fuego.Server, handler controller.APIResources) {

	event := fuego.Group(
		server,
		"/events",
		option.Summary("Events routes"),
		option.Description("Default description for all Events routes"),
		option.Tags("events"),
	)

	fuego.Get(event, "/", handler.GetAllEvents,
		option.Summary("Get all events"),
		option.Description("Returns a list of all events"),
	)

	fuego.Get(event, "/{id}", handler.GetEventByID,
		option.Summary("Get user by ID"),
		option.Description("Returns a single user by their ID"),
	)

	fuego.Post(event, "/", handler.CreateEvent,
		option.DefaultStatusCode(201),
		option.Summary("Create a new event"),
		option.Description("Create a single user using given informations"),
	)

	fuego.Put(event, "/{id}", handler.UpdateEvent,
		option.Summary("Update an event"),
		option.Description("Update an user by his ID using the given information"),
	)

	fuego.Delete(event, "/{id}", handler.DeleteEvent,
		option.DefaultStatusCode(204),
		option.Summary("Delete an user"),
		option.Description("Delete an user by his ID"),
	)
}
