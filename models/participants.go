package models

type Participants struct {
	IDParticipants int          `json:"id_participants" gorm:"primaryKey"`
	UserID         int          `json:"user_id" gorm:"not null"`
	User           User         `gorm:"foreignKey:UserID"`
	ConversationID int          `json:"conversation_id" gorm:"not null"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID"`
}

type ParticipantsRequest struct {
	UserID int
}
