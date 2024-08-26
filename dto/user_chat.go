package dto

type UserChatRequest struct {
	UserId    string `json:"userId"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
}

type UserChat struct {
	UserId       string `json:"userId"`
	Username     string `json:"username"`
	UserImageUrl string `json:"userImageUrl"`
	LastMessage  string `json:"lastMessage"`
	CreatedAt    int64  `json:"createdAt"`
}
