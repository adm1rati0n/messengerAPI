package controllers

import (
	"github.com/gofiber/fiber/v2"
	"messengerAPI/initializers"
	"messengerAPI/models"
	"strconv"
	"time"
)

func GetMessages(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(models.UserResponse)
	var messages []models.Message
	if result := initializers.DB.Find(&messages, "conversation_id = $1 "+
		"and not exists(select null from deleted_messages where deleted_messages.message_id = messages.id_message "+
		"and deleted_messages.user_id = $2) and not exists(select null from deleted_messages where "+
		"deleted_messages.message_id = messages.id_message "+
		"AND deleted_messages.is_deleted_for_all = true) "+
		"and exists(select null from participants where conversation_id = $1 and user_id = $2)", id, user.IDUser); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	var messagesView []models.MessageResponse

	for _, element := range messages {
		var messageView models.MessageResponse
		var attachments []models.Attachment
		if result := initializers.DB.Find(&attachments, "message_id = $1", element.IDMessage); result.Error != nil {
			return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
		}
		attachmentsResponse := models.FilterAttachmentsRecord(&attachments)

		var sender models.User
		if result := initializers.DB.First(&sender, element.SenderID); result.Error != nil {
			return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
		}
		senderResponse := models.FilterUserRecord(&sender)
		messageView = *models.FilterMessageRecord(&element, senderResponse, attachmentsResponse)
		messagesView = append(messagesView, messageView)
	}
	return c.Status(fiber.StatusOK).JSON(messagesView)
}

func GetDialogMessages(c *fiber.Ctx) error {
	id := c.Params("id")
	user := c.Locals("user").(models.UserResponse)
	var messages []models.Message
	if result := initializers.DB.Find(&messages, "dialog_id = $1 "+
		"and not exists(select null from deleted_messages where deleted_messages.message_id = messages.id_message "+
		"and deleted_messages.user_id = $2) and not exists(select null from deleted_messages where "+
		"deleted_messages.message_id = messages.id_message "+
		"AND deleted_messages.is_deleted_for_all = true)", id, user.IDUser); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	var messagesView []models.MessageResponse

	for _, element := range messages {
		var messageView models.MessageResponse
		var attachments []models.Attachment
		if result := initializers.DB.Find(&attachments, "message_id = $1", element.IDMessage); result.Error != nil {
			return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
		}
		attachmentsResponse := models.FilterAttachmentsRecord(&attachments)

		var sender models.User
		if result := initializers.DB.First(&sender, element.SenderID); result.Error != nil {
			return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
		}
		senderResponse := models.FilterUserRecord(&sender)
		messageView = *models.FilterMessageRecord(&element, senderResponse, attachmentsResponse)
		messagesView = append(messagesView, messageView)
	}
	return c.Status(fiber.StatusOK).JSON(messagesView)
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":        "bad request",
			"golangError:": err.Error(),
		})
	}
	var message models.Message
	message.ConversationID = &id
	message.SenderID = user.IDUser
	message.Text = &body.Text
	message.DateTime = time.Now()
	initializers.DB.Create(&message)
	return c.Status(fiber.StatusOK).JSON(&message)
}

func CreateDialogMessage(c *fiber.Ctx) error {
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
	message.DialogID = &id
	message.SenderID = user.IDUser
	message.Text = &body.Text
	message.DateTime = time.Now()
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
