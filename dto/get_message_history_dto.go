package dto

type GetMessageHistoryDTO struct {
	UserIdFrom string `json:"userIdFrom"`
	UserIdTo   string `json:"userIdTo"`
}
