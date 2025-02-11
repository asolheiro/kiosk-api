package controller

import "github.com/go-fuego/fuego"

func HealthCheck(c fuego.ContextNoBody) (string, error) {
	return "Kiosk API - v2.0.\nStatus: Healthy.\n", nil
}