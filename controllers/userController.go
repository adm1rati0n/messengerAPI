package controllers

import (
	"github.com/gofiber/fiber/v2"
	"messengerAPI/initializers"
	"messengerAPI/models"
)

func GetLoggedUser(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	if result := initializers.DB.Find(&users); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if result := initializers.DB.First(&user, id); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&user)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	body := models.UpdateUserRequest{}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var user models.User
	if result := initializers.DB.First(&user, id); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	user.Surname = body.Surname
	user.Name = body.Name
	user.MiddleName = body.MiddleName
	user.AvatarURL = body.AvatarURL
	user.Department = body.Department

	initializers.DB.Save(&user)
	return c.Status(fiber.StatusOK).JSON(&user)
}

func ChangePassword(c *fiber.Ctx) error {
	id := c.Params("id")
	body := models.UpdatePasswordRequest{}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var user models.User
	if result := initializers.DB.First(&user, id); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	if user.Password != body.OldPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Старый пароль не совпадает")
	} else if user.Password == body.NewPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Новый пароль совпадает с текущим")
	} else if body.NewPassword != body.ConfirmPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Пароли не совпадают")
	}
	user.Password = body.NewPassword

	initializers.DB.Save(&user)
	return c.Status(fiber.StatusOK).JSON(&user)
}

func ChangeDecryptPassword(c *fiber.Ctx) error {
	id := c.Params("id")
	body := models.UpdatePasswordRequest{}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var user models.User
	if result := initializers.DB.First(&user, id); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	if user.DecryptPassword != body.OldPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Старый пароль не совпадает")
	} else if user.DecryptPassword == body.NewPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Новый пароль совпадает с текущим")
	} else if body.NewPassword != body.ConfirmPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Пароли не совпадают")
	}
	user.DecryptPassword = body.NewPassword

	initializers.DB.Save(&user)
	return c.Status(fiber.StatusOK).JSON(&user)
}
