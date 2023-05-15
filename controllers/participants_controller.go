package controllers

import (
	"github.com/gofiber/fiber/v2"
	"messengerAPI/initializers"
	"messengerAPI/models"
	"strconv"
)

func AddParticipant(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Не смог считать айди чата",
		})
	}
	body := models.ParticipantsRequest{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Не смог считать тело",
		})
	}
	var participant models.Participants
	participant.UserID = body.UserID
	participant.ConversationID = id
	if result := initializers.DB.Create(&participant); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(&participant)
}

func GetConversationParticipants(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var participants []models.Participants
	if result := initializers.DB.Find(&participants, "conversation_id = $1", id); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&participants)
}

func DeleteParticipant(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var body int
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var participant models.User

	if result := initializers.DB.First(&participant, "conversation_id = $1 and user_id = $2", id, body); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	if result := initializers.DB.Delete(&participant); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}
