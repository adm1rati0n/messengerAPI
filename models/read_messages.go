package models

type ReadMessages struct {
	IDReadMessage int     `json:"id_read_message" gorm:"primaryKey"`
	UserID        int     `json:"user_id" gorm:"not null"`
	User          User    `gorm:"foreignKey:UserID"`
	MessageID     int     `json:"message_id" gorm:"not null"`
	Message       Message `gorm:"foreignKey:MessageID"`
}
