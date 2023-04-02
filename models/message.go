package models

import (
	"time"
)

type Message struct {
	IDMessage      int          `json:"id_message" gorm:"primaryKey"`
	Text           *string      `json:"text" gorm:"type:varchar(255)"`
	DateTime       time.Time    `json:"date_time" gorm:"type:timestamp;not null;default:now()"`
	SenderID       int          `json:"sender_id" gorm:"not null"`
	Sender         User         `gorm:"foreignKey:SenderID"`
	ConversationID *int         `json:"conversation_id"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID"`
	DialogID       *int         `json:"dialog_id"`
	Dialog         Dialog       `gorm:"foreignKey:DialogID"`
}

type MessageRequest struct {
	Text string
}

type MessageResponse struct {
	IDMessage   int                   `json:"id_message,omitempty"`
	Text        *string               `json:"text,omitempty"`
	DateTime    string                `json:"date_time"`
	Sender      UserResponse          `json:"sender"`
	Attachments *[]AttachmentResponse `json:"attachments"`
}

func FilterMessageRecord(message *Message, user UserResponse, attachments *[]AttachmentResponse) *MessageResponse {
	if message.IDMessage != 0 {
		return &MessageResponse{
			IDMessage:   message.IDMessage,
			Text:        message.Text,
			DateTime:    message.DateTime.Format("02.01.2006 15:04:05"),
			Sender:      user,
			Attachments: attachments,
		}
	} else {
		return nil
	}
}
