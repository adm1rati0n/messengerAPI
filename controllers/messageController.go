package controllers

import (
	"github.com/gofiber/fiber/v2"
	"messengerAPI/initializers"
	"messengerAPI/models"
	"strconv"
)

func GetMessages(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(models.UserResponse)
	var messages []models.Message
	if result := initializers.DB.Find(&messages, "conversation_id = $1 "+
		"and not exists(select null from deleted_messages where deleted_messages.message_id = messages.id_message "+
		"and deleted_messages.user_id = $2) and not exists(select null from deleted_messages where "+
		"deleted_messages.message_id = messages.id_message "+
		"AND deleted_messages.is_deleted_for_all = true)", id, user.IDUser); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	return c.Status(fiber.StatusOK).JSON(&messages)
}

func CreateMessage(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	body := models.MessageRequest{}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var message models.Message
	message.ConversationID = id
	message.SenderID = user.IDUser
	message.Text = body.Text
	initializers.DB.Create(&message)
	return c.Status(fiber.StatusOK).JSON(&message)
}

func AttachFile(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	body := models.AttachmentRequest{}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var attachment models.Attachment
	attachment.FileURL = body.FileURL
	attachment.MessageID = id
	initializers.DB.Create(&attachment)
	return c.Status(fiber.StatusOK).JSON(&attachment)
}

func DeleteMessage(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var message models.DeletedMessages
	message.IsDeletedForAll = false
	message.MessageID = id
	message.UserID = user.IDUser
	initializers.DB.Create(&message)
	return c.Status(fiber.StatusOK).JSON(&message)
}

func DeleteMessageForAll(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var message models.DeletedMessages
	message.IsDeletedForAll = true
	message.MessageID = id
	message.UserID = user.IDUser
	initializers.DB.Create(&message)
	return c.Status(fiber.StatusOK).JSON(&message)
}