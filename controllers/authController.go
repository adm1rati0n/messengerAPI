package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"messengerAPI/initializers"
	"messengerAPI/models"
	"strings"
	"time"
)

func SignUpUser(c *fiber.Ctx) error {
	body := models.SignUpRequest{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail1", "message": err.Error()})
	}
	errors := models.ValidateStruct(body)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})

	}
	//errors := models.ValidateStruct(payload)
	//if errors != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	//
	//}

	if body.Password != body.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail2", "message": "Passwords do not match"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail3", "message": err.Error()})
	}

	if body.DecryptPassword != body.ConfirmDecryptPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail4", "message": "Passwords do not match"})
	}

	hashedDecryptPassword, err := bcrypt.GenerateFromPassword([]byte(body.DecryptPassword), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail5", "message": err.Error()})
	}

	newUser := models.User{
		Surname:         body.Surname,
		Name:            body.Name,
		MiddleName:      &body.MiddleName,
		Login:           body.Login,
		Password:        string(hashedPassword),
		DecryptPassword: string(hashedDecryptPassword),
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail6", "message": "User with that login already exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error7", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": models.FilterUserRecord(&newUser)}})
}

func SignInUser(c *fiber.Ctx) error {
	body := models.AuthRequest{}
	fmt.Println("Метод вызван")
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	var user models.User
	result := initializers.DB.First(&user, "login = ?", body.Login)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad credentials",
		})
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad credentials",
		})
	}
	return GenerateJWTToken(c, user)
}
func LogoutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

//func (h handler) Login(c *fiber.Ctx) error {
//	body := models.AuthRequest{}
//	if err := c.BodyParser(&body); err != nil {
//		return fiber.NewError(fiber.StatusBadRequest, err.Error())
//	}
//	var user models.User
//	result := h.DB.First(&user, "Login = $1 and Password = $2", body.Login, body.Password)
//	if result.Error != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//			"error": "Bad credentials",
//		})
//	}
//	return GenerateJWTToken(c, user)
//}
