package routers

import (
	"github.com/asolheiro/kiosk-api/internal-v2/controller"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)



func UsersRouter(server *fuego.Server, handler controller.APIResources) {

	users := fuego.Group(server, "/users",
		option.Summary("Users routes"),
		option.Description("Default description for all Users routes"),
		option.Tags("users"),

	)

	fuego.Get(users, "/", handler.GetAllUsers)
	fuego.Get(users, "/{id}", handler.GetUserByID)
	fuego.Post(users, "/", handler.CreateUser, option.DefaultStatusCode(201))
	fuego.Put(users, "/{id}", handler.UpdateUser)
	fuego.Delete(users, "/{id}", handler.DeleteUser, option.DefaultStatusCode(204))
}