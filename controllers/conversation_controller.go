package controllers

import (
	"github.com/gofiber/fiber/v2"
	"messengerAPI/initializers"
	"messengerAPI/models"
	"strconv"
)

func CreateConversation(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)
	body := models.ConversationRequest{}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var conversation models.Conversation
	conversation.ConversationAvatarURL = body.ConversationAvatarURL
	conversation.Title = body.Title
	conversation.CreatorID = &user.IDUser
	initializers.DB.Save(&conversation)
	return c.Status(fiber.StatusOK).JSON(&conversation)
}

func GetAttachments(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var attachments []models.Attachment
	err = initializers.DB.Joins("join messages on messages.id_message = attachments.message_id").
		Joins("join conversations on conversations.id_conversation = messages.conversation_id").
		Where("conversations.id_conversation = ?", id).Find(&attachments).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving attachments",
		})
	}
	return c.JSON(attachments)
}

func CountUnreadMessages(userID int, conversationID int) int {
	var count int64
	initializers.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_id != ?", conversationID, userID).
		Where("date_time > (?)", initializers.DB.Table("read_messages").
			Select("MAX(messages.date_time)").
			Joins("INNER JOIN messages ON read_messages.message_id = messages.id_message").
			Where("read_messages.user_id = ? AND messages.conversation_id = ?", userID, conversationID)).
		Count(&count)

	return int(count)
}

func GetConversations(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)
	var conversations []models.Conversation

	if result := initializers.DB.Joins("LEFT JOIN messages ON "+
		"conversations.id_conversation = messages.conversation_id").
		Group("conversations.id_conversation").
		Order("MAX(messages.date_time) DESC").
		Find(&conversations, "exists(select null from participants where conversation_id = "+
			"conversations.id_conversation and user_id = ?)", user.IDUser); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	var conversationsView []models.ConversationWithLastMessage
	for _, element := range conversations {
		var message models.Message
		initializers.DB.Last(&message, "conversation_id = $1 and "+
			"not exists(select null from deleted_messages where "+
			"deleted_messages.message_id = messages.id_message and "+
			"deleted_messages.user_id = $2)", element.IDConversation, user.IDUser)

		var sender models.User
		initializers.DB.First(&sender, "id_user = ?", message.SenderID)

		var attachments []models.Attachment
		if result := initializers.DB.Find(&attachments, "message_id = $1", message.IDMessage); result.Error != nil {
			return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
		}
		conversationsView = append(conversationsView, models.FilterConversationRecord(&element, models.FilterSenderRecord(&sender), message, attachments, CountUnreadMessages(user.IDUser, element.IDConversation)))
	}

	return c.Status(fiber.StatusOK).JSON(&conversationsView)
}

func EditConversation(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	body := models.ConversationRequest{}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var conversation models.Conversation
	if result := initializers.DB.First(&conversation, id); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	conversation.Title = body.Title
	conversation.ConversationAvatarURL = body.ConversationAvatarURL
	initializers.DB.Save(&conversation)
	return c.Status(fiber.StatusOK).JSON(&conversation)
}
