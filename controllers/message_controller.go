package controllers

import (
	"github.com/gofiber/fiber/v2"
	"messengerAPI/initializers"
	"messengerAPI/models"
	"strconv"
	"time"
)

func GetMessages(c *fiber.Ctx) error {
	//var page = c.Query("page", "")
	//var limit = c.Query("limit", "30")
	//
	//intPage, _ := strconv.Atoi(page)
	//intLimit, _ := strconv.Atoi(limit)
	//offset := (intPage - 1) * intLimit

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
		senderResponse := models.FilterSenderRecord(&sender)
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
			"error": err.Error(),
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

func ReadMessage(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var message models.ReadMessages
	err = initializers.DB.Model(&message).Where("user_id = ? "+
		"and message_id = ?", user.IDUser, id).First(&message).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Что-то пошло не так",
		})
	}

	if message.IDReadMessage == 0 {
		var readMessage models.ReadMessages
		readMessage.MessageID = id
		readMessage.UserID = user.IDUser
		initializers.DB.Create(&message)
		return c.Status(fiber.StatusOK).JSON("Сообщение прочитано (create)")
	}
	message.UserID = user.IDUser
	message.MessageID = id
	initializers.DB.Save(&message)
	return c.Status(fiber.StatusOK).JSON("Сообщение прочитано (save)")
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
	message.MessageID = id
	message.UserID = user.IDUser
	initializers.DB.Create(&message)
	return c.Status(fiber.StatusOK).JSON(&message)
}

func DeleteMessageForAll(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var message models.Message
	err = initializers.DB.Model(&message).Where("id_message = ?", id).First(&message).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Что-то пошло не так",
		})
	}

	if message.IDMessage == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Сообщения с таким id не существует",
		})
	}

	err = initializers.DB.Delete(&message).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON("Сообщение удалено")
}
