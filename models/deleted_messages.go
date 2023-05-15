package models

type DeletedMessages struct {
	IDDeletedMessage int     `json:"id_deleted_message" gorm:"primaryKey"`
	UserID           int     `json:"user_id" gorm:"not null"`
	User             User    `gorm:"foreignKey:UserID"`
	MessageID        int     `json:"message_id" gorm:"not null"`
	Message          Message `gorm:"foreignKey:MessageID"`
}
