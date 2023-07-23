package dto

type GetMessageHistoryDTO struct {
	UserIdFrom string `json:"userIdFrom"`
	UserIdTo   string `json:"userIdTo"`
	PageIndex  int    `form:"pageIndex"`
	PageSize   int    `form:"pageSize"`
}
