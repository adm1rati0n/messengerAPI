package models

type Dialog struct {
	IDDialog   int  `json:"id_dialog" gorm:"primaryKey"`
	CreatorID  int  `json:"creator_id"`
	Creator    User `json:"creator" gorm:"foreignKey:CreatorID"`
	ReceiverID int  `json:"receiver_id"`
	Receiver   User `json:"receiver" gorm:"foreignKey:ReceiverID"`
}

type DialogResponse struct {
	IDDialog    int              `json:"id_dialog"`
	Companion   UserResponse     `json:"creator"`
	LastMessage *MessageResponse `json:"last_message"`
}

type DialogHelper struct {
	CreatorID  int
	ReceiverID int
}

func FilterDialogRecord(dialog *Dialog, user UserResponse, message Message, attachments []Attachment, companion UserResponse) DialogResponse {
	return DialogResponse{
		IDDialog:    dialog.IDDialog,
		Companion:   companion,
		LastMessage: FilterMessageRecord(&message, user, FilterAttachmentsRecord(&attachments)),
	}
}
