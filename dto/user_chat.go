package dto

type UserChatRequest struct {
	UserId    string `json:"userId"`
	pageIndex int    `json:"pageIndex"`
	pageSize  int64  `json:"pageSize"`
}

type UserChat struct {
	UserId       string `json:"userId"`
	Username     string `json:"username"`
	UserImageUrl string `json:"userImageUrl"`
	LastMessage  string `json:"lastMessage"`
	CreatedAt    int64  `json:"createdAt"`
}
