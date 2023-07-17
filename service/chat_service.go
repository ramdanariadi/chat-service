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
	GetUserWithLastMessage(requestBody dto.UserChatRequest) []*dto.UserChat
}

type ChatServiceImpl struct {
	*gorm.DB
}

func (service *ChatServiceImpl) GetUserWithLastMessage(requestBody dto.UserChatRequest) []*dto.UserChat {
	var userAndLastMessage []*dto.UserChat
	sql := "WITH ranked_message AS (SELECT *, RANK() OVER (PARTITION BY object ORDER BY created_at DESC ) AS rank_val " +
		"FROM (SELECT recipient AS object, message, created_at " +
		"	FROM chats " +
		"	WHERE sender = ? " +
		"	UNION " +
		"	SELECT sender AS object, message, created_at " +
		"	FROM chats " +
		"	WHERE recipient = ?) AS raw) " +
		"SELECT object AS UserId, message AS LastMessage " +
		"FROM ranked_message " +
		"WHERE rank_val = 1 " +
		"ORDER BY created_at DESC "
	log.Println("user id : " + requestBody.UserId)
	rows, err := service.DB.Raw(sql, requestBody.UserId, requestBody.UserId).Rows()
	utils.LogIfError(err)
	for rows.Next() {
		var userChat dto.UserChat
		err := rows.Scan(&userChat.UserId, &userChat.LastMessage)
		utils.LogIfError(err)
		userAndLastMessage = append(userAndLastMessage, &userChat)
	}
	return userAndLastMessage
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
