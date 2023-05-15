package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"messengerAPI/initializers"
	"messengerAPI/models"
	"time"
)

func GenerateJWTToken(c *fiber.Ctx, user models.User) error {
	config, _ := initializers.LoadConfig(".")
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	claims := tokenByte.Claims.(jwt.MapClaims)
	claims["sub"] = user.IDUser
	claims["exp"] = time.Now().UTC().Add(config.JwtExpiresIn).Unix()
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = time.Now().Unix()

	tokenString, err := tokenByte.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Generating jwt token failed",
			"message": err,
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   config.JwtMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": tokenString,
		"user":  user,
	})
}
