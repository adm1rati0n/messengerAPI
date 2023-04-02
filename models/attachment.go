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

type AttachmentResponse struct {
	IDAttachment int    `json:"id_attachment"`
	FileURL      string `json:"file_url"`
}

func FilterAttachmentsRecord(attachments *[]Attachment) *[]AttachmentResponse {
	var attachmentsResponse []AttachmentResponse
	for _, element := range *attachments {
		var attachmentResponse AttachmentResponse
		attachmentResponse.IDAttachment = element.IDAttachment
		attachmentResponse.FileURL = element.FileURL
		attachmentsResponse = append(attachmentsResponse, attachmentResponse)
	}
	return &attachmentsResponse
}
