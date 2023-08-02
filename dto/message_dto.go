package dto

type MessageDTO struct {
	Sender            string `json:"sender"`
	SenderName        string `json:"SenderName"`
	SenderImageUrl    string `json:"SenderImageUrl"`
	Recipient         string `json:"recipient"`
	RecipientName     string `json:"recipientName"`
	RecipientImageUrl string `json:"recipientImageUrl"`
	Message           string `json:"message"`
}
