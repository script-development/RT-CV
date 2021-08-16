package controller

import (
	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	v1 := app.Group(`/v1`, InsertData())

	authRoutes(v1)
	scraperRoutes(v1)
	controllerRoutes(v1)
}

type ProfilesCtx uint8
type AuthCtx uint8
type KeyCtx uint8
type LoggerCtx uint8

const (
	ProfilesCtxKey = ProfilesCtx(0)
	AuthCtxKey     = AuthCtx(0)
	KeyCtxKey      = KeyCtx(0)
	LoggerCtxKey   = LoggerCtx(0)
)
