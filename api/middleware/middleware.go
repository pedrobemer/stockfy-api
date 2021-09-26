package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
)

type Firebase struct {
	// Mandatory
	FirebaseAuth *auth.Client

	Authorizer     func(string, string) (*auth.Token, error)
	SuccessHandler fiber.Handler
	ErrorHandler   fiber.ErrorHandler
	ContextKey     string
}

func configFirebase(config Firebase) Firebase {
	cfg := config

	if cfg.ContextKey == "" {
		config.ContextKey = "user"
	}

	if cfg.FirebaseAuth == nil {
		panic("Please pass Firebase App in config")
	}

	// Default Authorizer function
	cfg.Authorizer = func(IDToken string, CurrentURL string) (*auth.Token, error) {
		// if cfg.FirebaseApp == nil {
		if cfg.FirebaseAuth == nil {
			return nil, errors.New("Missing Firebase App Object")
		}

		// Verify IDToken
		token, err := cfg.FirebaseAuth.VerifyIDToken(context.Background(),
			IDToken)

		// Throw error for bad token
		if err != nil {
			return nil, errors.New("Malformed Token")
		}

		if !token.Claims["email_verified"].(bool) {
			return nil, errors.New("Email not verified")
		}

		return token, nil
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

			if err.Error() == "Missing Firebase App Object" {
				return c.Status(fiber.StatusBadRequest).SendString("Missing or Invalid Firebase App Object")
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

func NewFirebase(config Firebase) fiber.Handler {
	var idToken string

	cfg := configFirebase(config)

	return func(c *fiber.Ctx) error {
		url := c.Method() + "::" + c.Path()

		authString := c.Get(fiber.HeaderAuthorization)
		if authString != "" {
			idToken = strings.Split(authString, "Bearer ")[1]
		}

		if len(idToken) == 0 || authString == "" {
			return cfg.ErrorHandler(c, errors.New("Missing Token"))
		}

		token, err := cfg.Authorizer(idToken, url)
		if err != nil {
			return cfg.ErrorHandler(c, err)
		}

		if token != nil {

			type user struct {
				emailVerified bool
				userID, email string
			}

			// Set authenticated user data into local context
			c.Locals(cfg.ContextKey, user{
				email:         token.Claims["email"].(string),
				emailVerified: token.Claims["email_verified"].(bool),
				userID:        token.Claims["user_id"].(string),
			})
			fmt.Println(c)

			return cfg.SuccessHandler(c)
		}
		return cfg.ErrorHandler(c, err)
	}
}
