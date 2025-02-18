package controller

import (
	"database/sql"
	"errors"
	"fmt"


	dto "github.com/asolheiro/kiosk-api/internal-v2/dtos"
	"github.com/asolheiro/kiosk-api/internal-v2/sqlitestore"
	"github.com/go-fuego/fuego"
)

type EventsResourcesInterface interface {
	GetEventByID(string, fuego.ContextNoBody) (*dto.Event, error)
	GetAllEvents(fuego.ContextNoBody) []dto.Event
	CreateEvent(fuego.ContextWithBody[sqlitestore.CreateEventParams]) (dto.Event, error)
	UpdateEvent(fuego.ContextWithBody[sqlitestore.UpdateUserParams]) (dto.Event, error)
	DeleteEvent(string, fuego.ContextNoBody) (any, error)
}

func (api *APIResources) GetEventByID(ctx fuego.ContextNoBody) (*dto.Event, error) {
	id := ctx.PathParam("id")
	event, err := api.repo.GetEvent(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fuego.NotFoundError{
				Title: "Event not found",
				Detail: "No event the provided ID was found.",
				Err: err,
			}
		}
		return nil, err
	}
	response := dto.FromSQLCEvent(event)
	return &response, nil
}

func (api *APIResources) GetAllEvents(ctx fuego.ContextNoBody) ([]dto.Event, error) {
	events, err := api.repo.ListEvents(ctx.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []dto.Event{}, nil
		}
		return nil, err
	}
	return dto.NewEventResponseList(events), nil
}

func (api *APIResources) CreateEvent(ctx fuego.ContextWithBody[sqlitestore.CreateEventParams]) (dto.Event, error) {
	input, err := ctx.Body()
	if err != nil {
		return dto.Event{}, fmt.Errorf("error decoding body.\nerr: %v", err)
	}

	event, err := api.repo.CreateEvent(ctx.Context(), input)
	if err != nil {
		return dto.Event{}, fmt.Errorf("error creating user. \nerr: %v", err)
	}

	return dto.FromSQLCEvent(event), nil
}

func (api *APIResources) UpdateEvent(ctx fuego.ContextWithBody[sqlitestore.UpdateEventParams]) (dto.Event, error) {
	input, err := ctx.Body()
	input.ID = ctx.PathParam("eventId")
	if err != nil {
		return dto.Event{}, fmt.Errorf("error decoding body.\nerr: %v", err)
	}

	event, err := api.repo.UpdateEvent(ctx.Context(), input)
	if err != nil {
		return dto.Event{}, fmt.Errorf("error updating event.\nerr: %v", err)
	}
	return dto.FromSQLCEvent(event), nil
}

func (api *APIResources) DeleteEvent(ctx fuego.ContextNoBody) (any, error) {
	eventId := ctx.PathParam("eventId")
	err := api.repo.SoftDeleteEvent(ctx.Context(), eventId)
	return nil, err
}