package models

type Conversation struct {
	IDConversation        int    `json:"id_conversation" gorm:"primaryKey"`
	ConversationAvatarURL string `json:"conversation_avatar_url" gorm:"type:varchar(100);not null;default:'default-chat.png"`
	Title                 string `json:"title" gorm:"type:varchar(100);not null"`
	CreatorID             *int   `json:"creator_id"`
	Creator               User   `gorm:"foreignKey:CreatorID"`
}

type ConversationRequest struct {
	ConversationAvatarURL string `json:"conversation_avatar_url"`
	Title                 string `json:"title"`
}

type ConversationWithLastMessage struct {
	IDConversation        int              `json:"id_conversation"`
	ConversationAvatarURL string           `json:"conversation_avatar_url"`
	Title                 string           `json:"title"`
	LastMessage           *MessageResponse `json:"last_message"`
	UnreadNumber          int              `json:"unread_number"`
}

func FilterConversationRecord(conversation *Conversation, user *UserResponse, message Message, attachments []Attachment, unreadNumber int) ConversationWithLastMessage {
	return ConversationWithLastMessage{
		IDConversation:        conversation.IDConversation,
		ConversationAvatarURL: conversation.ConversationAvatarURL,
		Title:                 conversation.Title,
		LastMessage:           FilterMessageRecord(&message, user, FilterAttachmentsRecord(&attachments)),
		UnreadNumber:          unreadNumber,
	}
}
