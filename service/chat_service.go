package service

import (
	"github.com/google/uuid"
	"github.com/ramdanariadi/chat-service/dto"
	"github.com/ramdanariadi/chat-service/model"
	"github.com/ramdanariadi/chat-service/utils"
	"gorm.io/gorm"
	"log"
)

type ChatService interface {
	GetMessageHistory(requestBody dto.GetMessageHistoryDTO) *[]dto.MessageHistoryDTO
	StoreMessage(message dto.MessageDTO)
}

type ChatServiceImpl struct {
	*gorm.DB
}

func (service *ChatServiceImpl) StoreMessage(message dto.MessageDTO) {
	newUUID, _ := uuid.NewUUID()
	chat := model.Chat{Id: newUUID.String(), Sender: message.Sender, Recipient: message.Recipient, Message: message.Message}
	service.DB.Save(&chat)
	log.Println("end")
}

func (service *ChatServiceImpl) GetMessageHistory(reqBody dto.GetMessageHistoryDTO) *[]dto.MessageHistoryDTO {
	chatHistory := make([]dto.MessageHistoryDTO, 0)
	var chats []*model.Chat

	senderAndRecipient := []string{reqBody.UserIdFrom, reqBody.UserIdTo}
	tx := service.DB.Model(&model.Chat{}).Where("sender IN ? AND recipient IN ?", senderAndRecipient, senderAndRecipient).Find(&chats)
	utils.LogIfError(tx.Error)

	for _, chat := range chats {
		dto := dto.MessageHistoryDTO{From: chat.Sender, To: chat.Recipient, Message: chat.Message, Time: chat.CreatedAt.Unix()}
		chatHistory = append(chatHistory, dto)
	}

	return &chatHistory
}
