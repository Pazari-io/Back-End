package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

// basic auth with secure password will do because url is not public
func Protected() fiber.Handler {

	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "password",
		},
		Realm: "Restricted",
	})
}
