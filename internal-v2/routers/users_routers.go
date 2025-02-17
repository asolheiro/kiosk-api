package routers

import (
	
	"github.com/asolheiro/kiosk-api/internal-v2/controller"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)


func UsersRouter(server *fuego.Server, handler controller.APIResources) {

	users := fuego.Group(
		server, 
		"/users",
		option.Summary("Users routes"),
		option.Description("Default description for all Users routes"),
		option.Tags("users"),
	)

    fuego.Get(users, "/", handler.GetAllUsers,
        option.Summary("Get all users"),
        option.Description("Returns a list of all users"),
	)

    fuego.Get(users, "/{id}", handler.GetUserByID,
        option.Summary("Get user by ID"),
        option.Description("Returns a single user by their ID"),
    )

	fuego.Post(users, "/", handler.CreateUser, 
		option.DefaultStatusCode(201),
		option.Summary("Create a new users"),
        option.Description("Create a single user using given informations"),
	)

	fuego.Put(users, "/{id}", handler.UpdateUser,
		option.Summary("Update an users"),
		option.Description("Update an user by his ID using the given information"),
	)

	fuego.Delete(users, "/{id}", handler.DeleteUser, 
		option.DefaultStatusCode(204),
		option.Summary("Delete an user"),
        option.Description("Delete an user by his ID"),
	)
}