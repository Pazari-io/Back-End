package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/keyauth/v2"
)


var (
	secret = "SKj2j2k32Ii32k42pPp24khg3odNDh2z"
	authList   = []string{secret}
	errMissing = &fiber.Error{
		Code:    403000,
		Message: "Missing API key",
	}
	errInvalid = &fiber.Error{
		Code:    403001,
		Message: "Invalid API key",
	}
)

func KeyProtected() fiber.Handler {
return (keyauth.New(keyauth.Config{
	
	Validator : validator,

  }))
}

func validator(ctx *fiber.Ctx, s string) (bool, error) {
	if s == "" {
		return false, errMissing
	}

	for _, val := range authList {
		if s == val {
			return true, nil
		}
	}

	return false, errInvalid
}