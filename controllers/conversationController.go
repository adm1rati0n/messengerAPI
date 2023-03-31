package controllers

import (
	"github.com/gofiber/fiber/v2"
	"messengerAPI/initializers"
	"messengerAPI/models"
	"strconv"
)

func CreateDialog(c *fiber.Ctx) error {
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
	var conversation models.Conversation
	conversation.ConversationAvatarURL = *user.AvatarURL
	conversation.Title = user.Surname + " " + user.Name + " " + *user.MiddleName
	if result := initializers.DB.Create(&conversation); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	var participant models.Participants
	participant.UserID = user.IDUser
	participant.ConversationID = conversation.IDConversation
	if result := initializers.DB.Create(&participant); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":         "failed to create a participant",
			"error_message": result.Error.Error(),
		})
	}
	var participant2 models.Participants
	loggedUser := c.Locals("user").(models.UserResponse)
	participant2.UserID = loggedUser.IDUser
	participant2.ConversationID = conversation.IDConversation
	if result := initializers.DB.Create(&participant2); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":         "failed to create a participant",
			"error_message": result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(&conversation)
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
	*conversation.CreatorID = user.IDUser
	initializers.DB.Save(&conversation)
	return c.Status(fiber.StatusOK).JSON(&conversation)
}

func GetConversations(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)
	var conversations []models.Conversation
	if result := initializers.DB.Find(&conversations); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}
	var message models.Message
	var conversationsView []models.ConversationWithLastMessage
	for _, element := range conversations {
		var conversationView models.ConversationWithLastMessage
		conversationView.IDConversation = element.IDConversation
		conversationView.ConversationAvatarURL = element.ConversationAvatarURL
		conversationView.Title = element.Title
		initializers.DB.Last(&message, "conversation_id = $1 and "+
			"not exists(select null from deleted_messages where "+
			"deleted_messages.message_id = messages.id_message and "+
			"deleted_messages.user_id = $2) and not exists(select null from deleted_messages "+
			"where deleted_messages.message_id = messages.id_message "+
			"AND deleted_messages.is_deleted_for_all = true)", element.IDConversation, user.IDUser)
		conversationView.LastMessage = message
		conversationsView = append(conversationsView, conversationView)
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
