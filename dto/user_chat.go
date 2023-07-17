package dto

type UserChatRequest struct {
	UserId    string `json:"userId"`
	pageIndex int    `json:"pageIndex"`
	pageSize  int64  `json:"pageSize"`
}

type UserChat struct {
	UserId      string `json:"userId"`
	LastMessage string `json:"lastMessage"`
}
