package models

import "time"

type Message struct {
	IDMessage      int          `json:"id_message" gorm:"primaryKey"`
	Text           string       `json:"text" gorm:"type:varchar(255)"`
	DateTime       time.Time    `json:"date_time" gorm:"type:timestamp;not null;default:now()"`
	SenderID       int          `json:"sender_id" gorm:"not null"`
	Sender         User         `gorm:"foreignKey:SenderID"`
	ConversationID int          `json:"conversation_id" gorm:"not null"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID"`
}

type MessageRequest struct {
	Text string
}
