package controller

import (
	"fmt"

	dto "github.com/asolheiro/kiosk-api/internal-v2/dtos"
	"github.com/asolheiro/kiosk-api/internal-v2/sqlitestore"
	"github.com/go-fuego/fuego"
)

type CheckInResourcesInterface interface {
	CreateCheckIn(fuego.ContextWithBody[sqlitestore.CreateCheckInParams]) (dto.CheckIn, error)
	FindCheckInById(fuego.ContextNoBody) (dto.CheckIn, error)
	ListCheckIns(fuego.ContextNoBody) ([]dto.CheckIn)
}

func (api *APIResources) CreateCheckIn(ctx fuego.ContextWithBody[sqlitestore.CreateCheckInParams]) (dto.CheckIn, error) {
	input, err := ctx.Body()
	if err != nil {
		return dto.CheckIn{}, fmt.Errorf("error decoding body. err: %w", err)
	}

	checkin, err := api.repo.CreateCheckIn(ctx.Context(), input)
	if err != nil {
		return dto.CheckIn{}, fmt.Errorf("error creating checkin. err: %w", err)
	}

	return dto.FromSQLCCheckIn(checkin), nil
}

func (api *APIResources) FindCheckInById(ctx fuego.ContextNoBody) (dto.CheckIn, error) {
	checkInId := ctx.PathParam("checkInId")
	checkIn, err := api.repo.GetCheckIn(ctx.Context(), checkInId)
	if err != nil {
		return dto.CheckIn{}, fmt.Errorf("error finding checkin. err: %w", err)
	}

	return dto.FromSQLCCheckIn(checkIn), nil
}

func (api *APIResources) ListCheckIns(ctx fuego.ContextNoBody) ([]dto.CheckIn, error) {
	checkIns, err := api.repo.ListCheckIns(ctx.Context())
	if err != nil {
		return nil, fmt.Errorf("error listing checkins. err: %w", err)
	}

	return dto.NewCheckInsList(checkIns), nil
}