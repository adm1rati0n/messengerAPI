package controllers

import (
	"github.com/gofiber/fiber/v2"
	"messengerAPI/initializers"
	"messengerAPI/models"
	"strconv"
)

func CreateDialog(c *fiber.Ctx) error {
	loggedUser := c.Locals("user").(models.UserResponse)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Нуб",
		})
	}
	var user models.User
	if result := initializers.DB.First(&user, id); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	var dialog models.Dialog
	dialog.CreatorID = loggedUser.IDUser
	dialog.ReceiverID = user.IDUser
	if result := initializers.DB.Create(&dialog); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(&dialog)
}

func GetDialogs(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)

	var dialogs []models.Dialog
	if result := initializers.DB.Find(&dialogs, "creator_id = $1 or receiver_id = $1", user.IDUser); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	var dialogsView []models.DialogResponse
	for _, element := range dialogs {
		var companionID int
		if element.ReceiverID == user.IDUser {
			companionID = element.CreatorID
		} else {
			companionID = element.ReceiverID
		}

		var companion models.User
		initializers.DB.First(&user, "id_user = $1", companionID)

		var message models.Message
		initializers.DB.Last(&message, "dialog_id = $1 and "+
			"not exists(select null from deleted_messages where "+
			"deleted_messages.message_id = messages.id_message and "+
			"deleted_messages.user_id = $2) and not exists(select null from deleted_messages "+
			"where deleted_messages.message_id = messages.id_message "+
			"AND deleted_messages.is_deleted_for_all = true)", element.IDDialog, user.IDUser,
		)
		var attachments []models.Attachment
		if result := initializers.DB.Find(&attachments, "message_id = $1", message.IDMessage); result.Error != nil {
			return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
		}
		dialogsView = append(dialogsView, models.FilterDialogRecord(&element, user, message, attachments, models.FilterUserRecord(&companion)))
	}
	return c.Status(fiber.StatusOK).JSON(&dialogsView)
}

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

func GetConversations(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)
	var conversations []models.Conversation

	if result := initializers.DB.Find(&conversations, "exists(select null from participants where conversation_id = conversations.id_conversation and user_id = $1)", user.IDUser); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	var conversationsView []models.ConversationWithLastMessage
	for _, element := range conversations {
		var message models.Message
		initializers.DB.Last(&message, "conversation_id = $1 and "+
			"not exists(select null from deleted_messages where "+
			"deleted_messages.message_id = messages.id_message and "+
			"deleted_messages.user_id = $2) and not exists(select null from deleted_messages "+
			"where deleted_messages.message_id = messages.id_message "+
			"AND deleted_messages.is_deleted_for_all = true)", element.IDConversation, user.IDUser)

		var attachments []models.Attachment
		if result := initializers.DB.Find(&attachments, "message_id = $1", message.IDMessage); result.Error != nil {
			return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
		}
		conversationsView = append(conversationsView, models.FilterConversationRecord(&element, user, message, attachments))
	}

	return c.Status(fiber.StatusOK).JSON(&conversationsView)
}

func EditConversation(c *fiber.Ctx) error {
	id := c.Params("id")
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
