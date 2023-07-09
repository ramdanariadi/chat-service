package dto

type MessageHistoryDTO struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
}
