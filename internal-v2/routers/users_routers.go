package routers

import (
	"errors"

	"github.com/asolheiro/kiosk-api/internal-v2/controller"
	dtos "github.com/asolheiro/kiosk-api/internal-v2/dtos"
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
        option.AddResponse(200, "Success", fuego.Response{
			Type: []dtos.UserResponse{}, 
			ContentTypes: []string{"application/json"},
		}),
		option.AddError(404, "Users not found", errors.New("new error")),
	)

    fuego.Get(users, "/{id}", handler.GetUserByID,
        option.Summary("Get user by ID"),
        option.Description("Returns a single user by their ID"),
        option.AddResponse(200, "Success", fuego.Response{
			Type: &dtos.UserResponse{}, 
			ContentTypes: []string{"application/json"},
		}),
        option.AddError(404, "User not found",  errors.New("new error")),
    )
	fuego.Post(users, "/", handler.CreateUser, 
		option.DefaultStatusCode(201),
		option.Summary("Create a new users"),
        option.Description("Create a single user using given informations"),
        option.AddResponse(200, "Success", fuego.Response{
			Type: dtos.UserResponse{}, 
			ContentTypes: []string{"application/json"},
		}),
		option.AddError(400, "Error creating user", fuego.Response{
			Type: struct{}{},
			ContentTypes: []string{"application/json"},
		}),
	)

	fuego.Put(users, "/{id}", handler.UpdateUser,
		option.Summary("Update an users"),
		option.Description("Update an user by his ID using the given information"),
		option.AddResponse(200, "Success", fuego.Response{
			Type: dtos.UserResponse{}, 
			ContentTypes: []string{"application/json"},
		}),
		option.AddError(400, "User not found", fuego.Response{
			Type: struct{}{},
			ContentTypes: []string{"application/json"},
		}),
	)

	fuego.Delete(users, "/{id}", handler.DeleteUser, 
		option.DefaultStatusCode(204),
		option.Summary("Delete an user"),
        option.Description("Delete an user by his ID"),
        option.AddResponse(200, "Success", fuego.Response{
			Type: dtos.UserResponse{}, 
			ContentTypes: []string{"application/json"},
		}),
		option.AddError(400, "User not found", fuego.Response{
			Type: struct{}{},
			ContentTypes: []string{"application/json"},
		}),
	)
}