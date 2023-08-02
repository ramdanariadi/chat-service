package service

import (
	"github.com/google/uuid"
	"github.com/ramdanariadi/chat-service/dto"
	"github.com/ramdanariadi/chat-service/model"
	"github.com/ramdanariadi/chat-service/utils"
	"gorm.io/gorm"
	"time"
)

type ChatService interface {
	GetMessageHistory(requestBody dto.GetMessageHistoryDTO) dto.MessageHistoryDTO
	StoreMessage(message dto.MessageDTO)
	GetUserWithLastMessage(requestBody dto.UserChatRequest) []*dto.UserChat
}

type ChatServiceImpl struct {
	*gorm.DB
}

func (service *ChatServiceImpl) GetUserWithLastMessage(requestBody dto.UserChatRequest) []*dto.UserChat {
	var userAndLastMessage []*dto.UserChat
	sql := "WITH ranked_message AS (SELECT *, RANK() OVER (PARTITION BY object ORDER BY created_at DESC ) AS rank_val " +
		"FROM ((SELECT recipient AS object, recipient_name AS object_name, recipient_image_url AS object_image_url, message, created_at " +
		"FROM chats " +
		"WHERE sender = ? ORDER BY created_at DESC) " +
		"UNION " +
		"(SELECT sender AS object, sender_name AS object_name, sender_image_url AS object_image_url, message, created_at " +
		"FROM chats " +
		"WHERE recipient = ? ORDER BY created_at DESC) " +
		") AS raw) " +
		"SELECT object AS user_id, object_name AS user_name, object_image_url AS user_image_url, message AS LastMessage, created_at " +
		"FROM ranked_message " +
		"WHERE rank_val = 1 " +
		"ORDER BY created_at DESC "
	rows, err := service.DB.Raw(sql, requestBody.UserId, requestBody.UserId).Rows()
	utils.LogIfError(err)
	for rows.Next() {
		var userChat dto.UserChat
		createdAt := time.Time{}
		err := rows.Scan(&userChat.UserId, &userChat.Username, &userChat.UserImageUrl, &userChat.LastMessage, &createdAt)
		userChat.CreatedAt = createdAt.Unix()
		utils.LogIfError(err)
		userAndLastMessage = append(userAndLastMessage, &userChat)
	}
	return userAndLastMessage
}

func (service *ChatServiceImpl) StoreMessage(message dto.MessageDTO) {
	newUUID, _ := uuid.NewUUID()
	chat := model.Chat{Id: newUUID.String(), Sender: message.Sender, SenderName: message.SenderName, SenderImageUrl: message.SenderImageUrl,
		Recipient: message.Recipient, RecipientName: message.RecipientName, RecipientImageUrl: message.RecipientImageUrl, Message: message.Message}
	service.DB.Save(&chat)
}

func (service *ChatServiceImpl) GetMessageHistory(reqBody dto.GetMessageHistoryDTO) dto.MessageHistoryDTO {
	chatHistory := make([]*dto.MessageHistoryItemDTO, 0)
	var chats []*model.Chat

	senderAndRecipient := []string{reqBody.UserIdFrom, reqBody.UserIdTo}
	tx := service.DB.Model(&model.Chat{}).Where("sender IN ? AND recipient IN ?", senderAndRecipient, senderAndRecipient)

	var count int64
	tx.Count(&count)

	tx.Order("created_at DESC").Limit(reqBody.PageSize).Offset(reqBody.PageSize * reqBody.PageIndex).Find(&chats)
	utils.LogIfError(tx.Error)

	for _, chat := range chats {
		dto := dto.MessageHistoryItemDTO{Id: chat.Id, From: chat.Sender, To: chat.Recipient, Message: chat.Message, Time: chat.CreatedAt.Unix()}
		chatHistory = append(chatHistory, &dto)
	}

	return dto.MessageHistoryDTO{Data: chatHistory, RecordsTotal: count, RecordsFiltered: len(chatHistory)}
}
