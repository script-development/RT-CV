package controller

import (
	"errors"
	"os"

	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// IMap is a wrapper around map[string]any that's faster to use
type IMap map[string]any

// Routes defines the routes used
func Routes(app *fiber.App, appVersion string, testing bool) {
	b := routeBuilder.New(app)

	b.Group(`/api/v1`, func(b *routeBuilder.Router) {
		b.Get(`/health`, getStatus(appVersion))

		b.Group(`/schema`, func(b *routeBuilder.Router) {
			b.Get(`/openAPI`, routeGetOpenAPISchema(b))
			b.Get(`/cv`, routeGetCvSchema)
		})

		b.Get(`/auth/keyinfo`, routeGetKeyInfo, requiresAuth(0))

		b.Group(`/scraper`, func(b *routeBuilder.Router) {
			b.Post(`/scanCV`, routeScraperScanCV)
			b.Group(`/scannedReferenceNrs`, func(b *routeBuilder.Router) {
				b.Get(``, scannedReferenceNrs)
				b.Group(`/since`, func(b *routeBuilder.Router) {
					b.Get(`/hours/:hours`, scannedReferenceNrs)
					b.Get(`/days/:days`, scannedReferenceNrs)
					b.Get(`/weeks/:weeks`, scannedReferenceNrs)
				})
			})
		}, requiresAuth(models.APIKeyRoleScraper|models.APIKeyRoleDashboard))

		secretsRoutes := func(b *routeBuilder.Router) {
			b.Get(``, routeGetSecrets)
			b.Group(`/:key`, func(b *routeBuilder.Router) {
				b.Delete(``, routeDeleteSecret)
				b.Put(``, routeUpdateOrCreateSecret)
				b.Get(`/:encryptionKey`, routeGetSecret)
			}, validKeyMiddleware())
		}
		b.Group(`/secrets`, func(b *routeBuilder.Router) {
			b.Group(`/myKey`,
				secretsRoutes,
				requiresAuth(models.APIKeyRoleAll),
				middlewareBindMyKey(),
			)
			b.Group(`/otherKey`, func(b *routeBuilder.Router) {
				// This route exposes a lot of user information that's why only the dashboard role can access it
				b.Get(``, routeGetAllSecretsFromAllKeys, requiresAuth(models.APIKeyRoleDashboard))
				b.Group(
					`/:keyID`,
					secretsRoutes,
					middlewareBindKey(),
					requiresAuth(models.APIKeyRoleInformationObtainer|models.APIKeyRoleDashboard),
				)
			})
		})
		b.Group(`/scraperUsers/:keyID`, func(b *routeBuilder.Router) {
			b.Get(``, routeGetScraperUsers)
			b.Patch(``, routePatchScraperUser, requiresAuth(models.APIKeyRoleAdmin|models.APIKeyRoleDashboard|models.APIKeyRoleController))
			b.Delete(``, routeDeleteScraperUser, requiresAuth(models.APIKeyRoleAdmin|models.APIKeyRoleDashboard|models.APIKeyRoleController))
		}, middlewareBindKey(), requiresAuth(models.APIKeyRoleAll))

		b.Group(`/profiles`, func(b *routeBuilder.Router) {
			// Profile routes that required the information obtainer role
			b.Group(``, func(b *routeBuilder.Router) {
				b.Get(`count`, routeGetProfilesCount)
				b.Get(``, routeAllProfiles)
				b.Post(`query`, routeQueryProfiles)
				b.Get(`/:profile`, routeGetProfile, middlewareBindProfile())
			}, requiresAuth(models.APIKeyRoleInformationObtainer|models.APIKeyRoleDashboard))

			// Profile routes that require the controller role
			b.Group(``, func(b *routeBuilder.Router) {
				b.Post(``, routeCreateProfile, requiresAuth(models.APIKeyRoleController))
				b.Group(`/:profile`, func(b *routeBuilder.Router) {
					b.Put(``, routeModifyProfile, requiresAuth(models.APIKeyRoleController))
					b.Delete(``, routeDeleteProfile, requiresAuth(models.APIKeyRoleController))
				}, middlewareBindProfile())
			}, requiresAuth(models.APIKeyRoleController|models.APIKeyRoleDashboard))
		}, requiresAuth(0))

		b.Group(`/keys`, func(b *routeBuilder.Router) {
			b.Get(``, routeGetKeys)
			b.Get(`/scrapers`, routeGetScraperKeys)
			b.Post(``, routeCreateKey)
			b.Group(`/:keyID`, func(b *routeBuilder.Router) {
				b.Get(``, routeGetKey)
				b.Put(``, routeUpdateKey)
				b.Delete(``, routeDeleteKey)
			}, middlewareBindKey())
		}, requiresAuth(models.APIKeyRoleDashboard))

		b.Group(`/onMatchHooks`, func(b *routeBuilder.Router) {
			b.Get(``, routeGetOnMatchHooks)
			b.Post(``, routeCreateOnMatchHooks)
			b.Group(`/:hookID`, func(b *routeBuilder.Router) {
				b.Delete(``, routeDeleteOnMatchHook)
				b.Put(``, routeUpdateOnMatchHook)
				b.Post(`/test`, routeTestOnMatchHook)
			}, middlewareBindHook())
		}, requiresAuth(models.APIKeyRoleController|models.APIKeyRoleInformationObtainer|models.APIKeyRoleDashboard))

		b.Group(`/matcherTree`, func(b *routeBuilder.Router) {
			b.Get(``, routeGetMatcherTree)
			b.Post(`/addLeaf`, routeAddMatcherLeaf)
			b.Group(`/:id`, func(b *routeBuilder.Router) {
				b.Get(``, routeGetMatcherTree)
				b.Put(``, routePutMatcherBranch)
				b.Delete(``, routeDeleteMatcherBranch)
				b.Post(`/addLeaf`, routeAddMatcherLeaf)
			})
		}, requiresAuth(models.APIKeyRoleController|models.APIKeyRoleInformationObtainer|models.APIKeyRoleDashboard))
	})

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
		staticOpts := fiber.Static{Compress: true}
		b.Static("", "./dashboard/out/index.html", staticOpts)
		b.Static("login", "./dashboard/out/login.html", staticOpts)
		b.Static("tryMatcher", "./dashboard/out/tryMatcher.html", staticOpts)
		b.Static("tryPdfGenerator", "./dashboard/out/tryPdfGenerator.html", staticOpts)
		b.Static("docs", "./dashboard/out/docs.html", staticOpts)
		b.Static("_next", "./dashboard/out/_next", staticOpts)
		b.Static("favicon.ico", "./dashboard/out/favicon.ico", staticOpts)
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
}

func getStatus(appVersion string) routeBuilder.R {
	return routeBuilder.R{
		Description: "Get the server status",
		Res:         GetStatusResponse{},
		Fn: func(c *fiber.Ctx) error {
			return c.JSON(GetStatusResponse{
				Status:     true,
				AppVersion: appVersion,
			})
		},
	}
}
