package controller

import (
	"fmt"

	imageservice "github.com/asolheiro/kiosk-api/internal-v2/imageService"
	"github.com/go-fuego/fuego"
)

type PostPrintParams struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	TemplateImage string `json:"template_image"`
}

func (api *APIResources) PostPrint(ctx fuego.ContextWithBody[PostPrintParams]) (any, error) {
	body, err := ctx.Body()
	if err != nil {
		return nil, fmt.Errorf("error decoding body, err: %w", err)
	}

	mainConfig, err := api.repo.GetFirstConfig(ctx.Context())
	if err != nil {
		return nil, fmt.Errorf("config not found, err: %w", err)
	}

	err = imageservice.DefaultPrint(
		mainConfig.TemplateImage.String,
		body.Title,
		body.Description,
		int(mainConfig.PositionX.Int64),
		int(mainConfig.PositionY.Int64),
		int(mainConfig.FontSize.Int64),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting image, err %w", err)
	}
	return nil, nil
}