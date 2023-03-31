package models

type Attachment struct {
	IDAttachment int     `json:"id_attachment" gorm:"primaryKey"`
	FileURL      string  `json:"file_url" gorm:"type:varchar(100);not null"`
	MessageID    int     `json:"message_id" gorm:"not null"`
	Message      Message `gorm:"foreignKey:MessageID"`
}

type AttachmentRequest struct {
	FileURL string `json:"file_url"`
}
