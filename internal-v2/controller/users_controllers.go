package controller

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/asolheiro/kiosk-api/internal-v2/sqlitestore"
	"github.com/go-fuego/fuego"
)

type UsersResourcesInterface interface {
	GetUserByID(string, fuego.ContextNoBody) (*sqlitestore.User, error)
	GetAllUsers(fuego.ContextNoBody) ([]sqlitestore.User)
	CreateUser(fuego.ContextWithBody[sqlitestore.CreateUserParams]) (sqlitestore.User, error)
	UpdateUser(fuego.ContextWithBody[sqlitestore.UpdateUserParams]) (sqlitestore.User, error)
	DeleteUser(string, fuego.ContextNoBody) (error)
}

func (api *APIResources) GetUserByID(ctx fuego.ContextNoBody) (*sqlitestore.User, error) {
	id := ctx.PathParam("id")
	var user sqlitestore.User
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
	return &user, nil
}

func (api *APIResources) GetAllUsers(ctx fuego.ContextNoBody) ([]sqlitestore.User, error) {
	var users []sqlitestore.User
	users, err := api.repo.ListUsers(ctx.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, nil
			}
		}

	return users, nil
}

func (api *APIResources) CreateUser(ctx fuego.ContextWithBody[sqlitestore.CreateUserParams]) (user sqlitestore.User, err error) {
	input, err := ctx.Body()
	if err != nil {
		return user, fmt.Errorf("error decoding body.\nerr: %v", err)
	}
	
	user, err = api.repo.CreateUser(ctx.Context(), input)
	if err != nil {
		return user, fmt.Errorf("error creating user. \nerr: %v", err)
	}

	return user, nil
}

func (api *APIResources) UpdateUser(ctx fuego.ContextWithBody[sqlitestore.UpdateUserParams]) (user sqlitestore.User, err error) {
	input, err := ctx.Body()
	if err != nil {
		return user, fmt.Errorf("error decoding body.\nerr: %v", err)
	}

	user, err = api.repo.UpdateUser(ctx.Context(), input)
	if err != nil {
		return user, fmt.Errorf("error updating user.\nerr: %v", err)
	}

	return user, nil 
}

func (api *APIResources) DeleteUser(ctx fuego.ContextNoBody) (any, error) {
	userId := ctx.PathParam("userId")
	err := api.repo.SoftDeleteUser(ctx.Context(), userId)
	return nil, err
}