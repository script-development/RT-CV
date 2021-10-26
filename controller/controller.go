package controller

import (
	"errors"
	"os"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// IMap is a wrapper around map[string]interface{} that's faster to use
type IMap map[string]interface{}

// Routes defines the routes used
func Routes(app *fiber.App, appVersion string, dbConn db.Connection, serverSeed string, testing bool) {
	b := routeBuilder.New(app)

	b.Group(`/api/v1`, func(b *routeBuilder.Router) {
		b.Get(`/health`, getStatus(appVersion))

		b.Group(`/schema`, func(b *routeBuilder.Router) {
			b.Get(`/openAPI`, routeGetOpenAPISchema(b))
			b.Get(`/cv`, routeGetCvSchema)
		})

		b.Group(`/auth`, func(b *routeBuilder.Router) {
			b.Get(`/keyinfo`, routeGetKeyInfo, requiresAuth(0))
			b.Get(`/seed`, routeAuthSeed)
		})

		b.Group(`/scraper`, func(b *routeBuilder.Router) {
			b.Post(`/scanCV`, routeScraperScanCV)
		}, requiresAuth(models.APIKeyRoleScraper|models.APIKeyRoleDashboard))

		secretsRoutes := func(b *routeBuilder.Router) {
			b.Get(``, routeGetSecrets)
			b.Group(`/:key`, func(b *routeBuilder.Router) {
				b.Delete(``, routeDeleteSecret)
				b.Put(``, routeUpdateOrCreateSecret)
				b.Get(`/:encryptionKey`, routeGetSecret)
			}, validKeyMiddleware())
		}
		b.Group(`/secrets/myKey`, secretsRoutes, requiresAuth(models.APIKeyRoleAll), middlewareBindMyKey())
		b.Group(`/secrets/otherKey`, func(b *routeBuilder.Router) {
			// This route exposes a lot of user information that's why only the dashboard role can access it
			b.Get(``, routeGetAllSecretsFromAllKeys, requiresAuth(models.APIKeyRoleDashboard))
			b.Group(
				`/:keyID`,
				secretsRoutes,
				middlewareBindKey(),
				requiresAuth(models.APIKeyRoleInformationObtainer|models.APIKeyRoleDashboard),
			)
		})

		b.Group(`/analytics`, func(b *routeBuilder.Router) {
			b.Group(`/matches`, func(b *routeBuilder.Router) {
				b.Get(`/period/:from/:to`, routeGetMatchesPeriod)
			})
		}, requiresAuth(models.APIKeyRoleInformationObtainer|models.APIKeyRoleDashboard))

		b.Group(`/control`, func(b *routeBuilder.Router) {
			b.Group(`/profiles`, func(b *routeBuilder.Router) {
				b.Post(``, routeCreateProfile)
				b.Get(``, routeAllProfiles)
				b.Group(`/:profile`, func(b *routeBuilder.Router) {
					b.Get(``, routeGetProfile)
					// b.Put(``, routeModifyProfile) // TODO
					b.Delete(``, routeDeleteProfile)
				}, middlewareBindProfile())
			})
		}, requiresAuth(models.APIKeyRoleController))

		b.Group(`/keys`, func(b *routeBuilder.Router) {
			b.Get(``, routeGetKeys)
			b.Post(``, routeCreateKey)
			b.Group(`/:keyID`, func(b *routeBuilder.Router) {
				b.Get(``, routeGetKey)
				b.Put(``, routeUpdateKey)
				b.Delete(``, routeDeleteKey)
			}, middlewareBindKey())
		}, requiresAuth(models.APIKeyRoleDashboard))
	}, InsertData(dbConn, serverSeed))

	_, err := os.Stat("./dashboard/out")
	if err == os.ErrNotExist {
		if !testing {
			log.Warn("dashboard not build, you won't be able to use the dashboard")
		}
	} else if err != nil {
		if !testing {
			log.WithError(err).Warn("unable to set dashboard routes, you won't be able to use the dashboard")
		}
	} else {
		// FIXME we currently need to manually add every dashboard route here.
		// It would be nice if these where auto generated
		b.Static("", "./dashboard/out/index.html", fiber.Static{Compress: true})
		b.Static("login", "./dashboard/out/login.html", fiber.Static{Compress: true})
		b.Static("tryMatcher", "./dashboard/out/tryMatcher.html", fiber.Static{Compress: true})
		b.Static("docs", "./dashboard/out/docs.html", fiber.Static{Compress: true})
		b.Static("_next", "./dashboard/out/_next", fiber.Static{Compress: true})
		b.Static("favicon.ico", "./dashboard/out/favicon.ico", fiber.Static{Compress: true})
		app.Use(func(c *fiber.Ctx) error {
			// 404 page
			return c.Status(404).SendFile("./dashboard/out/404.html", true)
		})
	}
}

// FiberErrorHandler handles errors in fiber
// In our case that means we change the errors from text to json
func FiberErrorHandler(c *fiber.Ctx, err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrorRes(c, 404, errors.New("item not found"))
	}
	return ErrorRes(c, 500, err)
}

// ErrorRes returns the error response
func ErrorRes(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(IMap{
		"error": err.Error(),
	})
}

// GetStatusResponse contains the response data for the getStatus route
type GetStatusResponse struct {
	Status     bool   `json:"status"`
	AppVersion string `json:"appVersion"`
	Seed       string `json:"seed"`
}

func getStatus(appVersion string) routeBuilder.R {
	return routeBuilder.R{
		Description: "Get the server status",
		Res:         GetStatusResponse{},
		Fn: func(c *fiber.Ctx) error {
			return c.JSON(GetStatusResponse{
				Status:     true,
				AppVersion: appVersion,
				Seed:       string(ctx.GetAuth(c).GetBaseSeed()),
			})
		},
	}
}
