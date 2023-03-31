package models

type Conversation struct {
	IDConversation        int    `json:"id_conversation" gorm:"primaryKey"`
	ConversationAvatarURL string `json:"conversation_avatar_url" gorm:"type:varchar(100);not null;default:'default.png"`
	Title                 string `json:"title" gorm:"type:varchar(100);not null"`
	CreatorID             *int   `json:"creator_id" `
	Creator               User   `gorm:"foreignKey:CreatorID"`
}

type ConversationRequest struct {
	ConversationAvatarURL string `json:"conversation_avatar_url"`
	Title                 string `json:"title"`
}

type ConversationWithLastMessage struct {
	IDConversation        int    `json:"id_conversation"`
	ConversationAvatarURL string `json:"conversation_avatar_url"`
	Title                 string `json:"title"`
	LastMessage           Message
}
