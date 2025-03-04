package controller

import (
	"database/sql"
	"errors"
	"fmt"

	dto "github.com/asolheiro/kiosk-api/internal-v2/dtos"
	"github.com/asolheiro/kiosk-api/internal-v2/sqlitestore"
	"github.com/go-fuego/fuego"
)

type GuestResourcesInterface interface {
	CreateUser(fuego.ContextWithBody[sqlitestore.CreateGuestParams]) (dto.Guest, error)
	ListUsers(fuego.ContextNoBody) ([]dto.Guest, error)
	FindGuestById(fuego.ContextNoBody) (dto.Guest, error)
	FindGuestByDocument(fuego.ContextNoBody) (dto.Guest, error)
	UpdateGuest(fuego.ContextWithBody[sqlitestore.UpdateGuestParams]) (dto.Guest, error)
	DeleteGuest(fuego.ContextNoBody) (any, error)
}

func (api *APIResources) CreateGuest(ctx fuego.ContextWithBody[sqlitestore.CreateGuestParams]) (dto.Guest, error) {
	input, err := ctx.Body()
	if err != nil {
		return dto.Guest{}, fmt.Errorf("error decoding body. err: %w", err)
	}

	guest, err := api.repo.CreateGuest(ctx.Context(), input)
	if err != nil {
		return dto.Guest{}, fmt.Errorf("error creating user, err: %w", err)
	}

	return dto.FromSQLGuest(guest), nil
}

func (api *APIResources) ListGuests(ctx fuego.ContextNoBody) ([]dto.Guest, error) {
	guests, err := api.repo.ListGuests(ctx.Context())
	if err != nil {
		return nil, fmt.Errorf("error finding guests, err: %w", err)
	}

	return dto.NewGuestsList(guests), nil
}

func (api *APIResources) FindGuestById(ctx fuego.ContextNoBody) (dto.Guest, error) {
	guestId := ctx.PathParam("guestId")
	guest, err := api.repo.GetGuest(ctx.Context(), guestId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.Guest{}, fuego.NotFoundError{
				Title: "Guest not found",
				Detail: "No guest with the provided ID was found",
				Err: err,
			}
		}
		return dto.Guest{}, err
	}
	return dto.FromSQLGuest(guest), nil
}

func (api *APIResources) FindGuestByDocument(ctx fuego.ContextNoBody) (dto.Guest, error) {
	guestDocument := ctx.PathParam("guestDocument")
	guest, err := api.repo.GetGuestByDocumentNumber(ctx.Context(), guestDocument)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.Guest{}, fuego.NotFoundError{
				Title: "Guest not found",
				Detail: "No guest with the provided document number was found",
				Err: err,
			}
		}
		return dto.Guest{}, fmt.Errorf("error finding guest, err: %w", err)
	}
	guestDto := dto.FromSQLGuest(guest)
	return guestDto, nil
}

func (api *APIResources) UpdateGuest(ctx fuego.ContextWithBody[sqlitestore.UpdateGuestParams]) (dto.Guest, error) {
	input, err := ctx.Body()
	input.ID = ctx.PathParam("guestId")
	if err != nil {
		return dto.Guest{}, fmt.Errorf("error decoding body, err: %w", err)
	}

	guest, err := api.repo.UpdateGuest(ctx.Context(), input)
	if err != nil {
		return dto.Guest{}, fmt.Errorf("error updating guest, err: %w", err)
	}
	return dto.FromSQLGuest(guest), nil
	
}

func (api *APIResources) DeleteGuest(ctx fuego.ContextNoBody) (any, error) {
	guestId := ctx.PathParam("guestId")
	err := api.repo.SoftDeleteGuest(ctx.Context(), guestId)
	if err != nil {
		return nil, fmt.Errorf("error deleting guest, err: %w", err) 
	}
	return nil, nil
}