package controller

import (
	"database/sql"
	"errors"
	"fmt"

	dto "github.com/asolheiro/kiosk-api/internal-v2/dtos"
	"github.com/asolheiro/kiosk-api/internal-v2/sqlitestore"
	"github.com/go-fuego/fuego"
)

type UsersResourcesInterface interface {
	GetUserByID(string, fuego.ContextNoBody) (*dto.User, error)
	GetAllUsers(fuego.ContextNoBody) ([]dto.User)
	CreateUser(fuego.ContextWithBody[sqlitestore.CreateUserParams]) (dto.User, error)
	UpdateUser(fuego.ContextWithBody[sqlitestore.UpdateUserParams]) (dto.User, error)
	DeleteUser(string, fuego.ContextNoBody) (any, error)
}

func (api *APIResources) GetUserByID(ctx fuego.ContextNoBody) (*dto.UserResponse, error) {
    id := ctx.PathParam("id")
    user, err := api.repo.GetUser(ctx.Context(), id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fuego.NotFoundError{
                Title:  "User not found",
                Detail: "No user with the provided ID was found.",
                Err:    err,
            }
        }
        return nil, err
    }
    response := dto.NewUserResponse(user)
    return &response, nil
}

func (api *APIResources) GetAllUsers(ctx fuego.ContextNoBody) ([]dto.UserResponse, error) {
    users, err := api.repo.ListUsers(ctx.Context())
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return []dto.UserResponse{}, nil
        }
        return nil, err
    }
    return dto.NewUserResponseList(users), nil
}

func (api *APIResources) CreateUser(ctx fuego.ContextWithBody[sqlitestore.CreateUserParams]) (userDto dto.User, err error) {
	input, err := ctx.Body()
	if err != nil {
		return dto.User{}, fmt.Errorf("error decoding body.\nerr: %v", err)
	}
	
	user, err := api.repo.CreateUser(ctx.Context(), input)
	if err != nil {
		return dto.User{}, fmt.Errorf("error creating user. \nerr: %v", err)
	}
	
	userDto = dto.FromSQLCUser(user)
	return userDto, nil
}

func (api *APIResources) UpdateUser(ctx fuego.ContextWithBody[sqlitestore.UpdateUserParams]) (userDto dto.User, err error) {
	input, err := ctx.Body()
	if err != nil {
		return dto.User{}, fmt.Errorf("error decoding body.\nerr: %v", err)
	}

	user, err := api.repo.UpdateUser(ctx.Context(), input)
	if err != nil {
		return dto.User{}, fmt.Errorf("error updating user.\nerr: %v", err)
	}
	userDto = dto.FromSQLCUser(user)
	return userDto, nil 
}

func (api *APIResources) DeleteUser(ctx fuego.ContextNoBody) (any, error) {
	userId := ctx.PathParam("userId")
	err := api.repo.SoftDeleteUser(ctx.Context(), userId)
	return nil, err
}