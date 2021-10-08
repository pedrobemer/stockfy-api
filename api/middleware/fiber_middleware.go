package middleware

import (
	"errors"
	"stockfyApi/entity"
	"stockfyApi/usecases/user"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type FiberMiddleware struct {
	// Mandatory
	UserAuthentication *user.Application

	Authorizer     func(string, string) (*entity.UserTokenInfo, error)
	SuccessHandler fiber.Handler
	ErrorHandler   fiber.ErrorHandler
	ContextKey     string
}

func configMiddleware(config FiberMiddleware) FiberMiddleware {
	cfg := config

	if cfg.ContextKey == "" {
		config.ContextKey = "user"
	}

	if cfg.UserAuthentication == nil {
		panic("Please pass use cases object for User package in config")
	}

	// Default Authorizer function
	cfg.Authorizer = func(IDToken string, CurrentURL string) (*entity.UserTokenInfo,
		error) {
		if cfg.UserAuthentication == nil {
			return nil, errors.New("Missing User package Object")
		}

		// Verify IDToken
		tokenInfo, err := cfg.UserAuthentication.UserTokenVerification(IDToken)

		// Throw error for bad token
		if err != nil {
			return nil, errors.New("Malformed Token")
		}

		if !tokenInfo.EmailVerified {
			return nil, errors.New("Email not verified")
		}

		return tokenInfo, nil
	}

	// Default Error Handler
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing Token" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or malformed Token")
			}

			if err.Error() == "Malformed Token" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or malformed Token")
			}

			if err.Error() == "Email not verified" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or malformed Token")
			}

			if err.Error() == "Missing User package Object" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or Invalid User package Object")
			}

			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Token")

		}
	}

	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	return cfg
}

func NewFiberMiddleware(config FiberMiddleware) fiber.Handler {
	var idToken string

	cfg := configMiddleware(config)

	return func(c *fiber.Ctx) error {
		url := c.Method() + "::" + c.Path()

		authString := c.Get(fiber.HeaderAuthorization)
		if authString != "" {
			idToken = strings.Split(authString, "Bearer ")[1]
		}

		if len(idToken) == 0 || authString == "" {
			return cfg.ErrorHandler(c, errors.New("Missing Token"))
		}

		tokenInfo, err := cfg.Authorizer(idToken, url)
		if err != nil {
			return cfg.ErrorHandler(c, err)
		}

		if tokenInfo != nil {

			type user struct {
				emailVerified bool
				userID, email string
			}

			// Set authenticated user data into local context
			c.Locals(cfg.ContextKey, user{
				email:         tokenInfo.Email,
				emailVerified: tokenInfo.EmailVerified,
				userID:        tokenInfo.UserID,
			})

			return cfg.SuccessHandler(c)
		}
		return cfg.ErrorHandler(c, err)
	}
}
