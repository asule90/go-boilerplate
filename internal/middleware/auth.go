package middleware

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/sule/go-boilerplate/config"
	"github.com/sule/go-boilerplate/pkg/errr"
	"google.golang.org/api/option"
)

const (
	contextKey    = "uid"
	jwtContextKey = "jwt_claims"
)

// FirebaseAuth holds the Firebase auth client.
type FirebaseAuth struct {
	client *firebaseAuth.Client
}

// InitFirebase initializes the Firebase app and returns a FirebaseAuth instance.
func InitFirebase(cfg *config.Provider) (*FirebaseAuth, error) {
	var app *firebase.App
	var err error

	if cfg.Auth.FirebaseCredentials != "" {
		opt := option.WithCredentialsFile(cfg.Auth.FirebaseCredentials)
		app, err = firebase.NewApp(context.Background(), nil, opt)
	} else {
		app, err = firebase.NewApp(context.Background(), nil)
	}
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	return &FirebaseAuth{client: client}, nil
}

// Auth returns a Fiber middleware that validates Firebase ID tokens.
func (fa *FirebaseAuth) Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractBearerToken(c)
		if token == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization token"})
		}

		decoded, err := fa.client.VerifyIDToken(context.Background(), token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals(contextKey, decoded.UID)
		return c.Next()
	}
}

// UID extracts the Firebase UID from the request context.
func UID(c *fiber.Ctx) string {
	uid, _ := c.Locals(contextKey).(string)
	return uid
}

// JWTAuth returns a Fiber middleware that validates JWT tokens.
func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := extractBearerToken(c)
		if token == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization token"})
		}

		claims, err := ParseJWTClaims(token, secret)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals(jwtContextKey, claims)
		return c.Next()
	}
}

// ParseJWTClaims parses and validates a JWT token string.
func ParseJWTClaims(tokenStr, secret string) (jwtv5.MapClaims, error) {
	token, err := jwtv5.Parse(tokenStr, func(token *jwtv5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errr.New(http.StatusUnauthorized, "unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwtv5.MapClaims)
	if !ok || !token.Valid {
		return nil, errr.New(http.StatusUnauthorized, "invalid token claims")
	}

	return claims, nil
}

// UserID extracts the subject claim (user ID) from a JWT in the request context.
func UserID(c *fiber.Ctx) string {
	claims, ok := c.Locals(jwtContextKey).(jwtv5.MapClaims)
	if !ok {
		return ""
	}
	sub, _ := claims["sub"].(string)
	return sub
}

func extractBearerToken(c *fiber.Ctx) string {
	header := c.Get("Authorization")
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}
