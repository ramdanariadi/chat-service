package service

import "github.com/ramdanariadi/chat-service/dto"

type ChatService interface {
	GetMessageHistory(requestBody dto.GetMessageHistoryDTO) dto.MessageHistoryDTO
	StoreMessage(message dto.MessageDTO)
	GetUserWithLastMessage(requestBody dto.UserChatRequest) []*dto.UserChat
}
