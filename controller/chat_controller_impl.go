package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ramdanariadi/chat-service/dto"
	"github.com/ramdanariadi/chat-service/service"
	"github.com/ramdanariadi/chat-service/utils"
	"gorm.io/gorm"
	"log"
)

type ChatControllerImpl struct {
	*gorm.DB
	Upgrader       *websocket.Upgrader
	Connections    []*service.UserConnection
	ConnectionChan chan *service.UserConnection
	Service        service.ChatService
}

func (contoller *ChatControllerImpl) WSHandler(ctx *gin.Context) {
	conn, err := contoller.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	utils.LogIfError(err)
	userId := ctx.Param("userId")
	connection := service.UserConnection{UserId: userId, Connection: conn}
	contoller.Connections = append(contoller.Connections, &connection)
	marshal, err := json.Marshal(dto.MessageDTO{Recipient: userId, Message: "hi!!"})
	utils.LogIfError(err)

	err = conn.WriteMessage(websocket.TextMessage, marshal)
	log.Println(err)

	go func(userConnection *service.UserConnection, userChan chan<- *service.UserConnection) {
		for {
			_, p, e := userConnection.Connection.ReadMessage()
			if nil != e {
				log.Println("here")
				utils.LogIfError(e)
				if websocket.IsUnexpectedCloseError(e) {
					log.Println("ws closed")
					userChan <- userConnection
				}
				break
			}
			log.Println(p)
			var message dto.MessageDTO
			err2 := json.Unmarshal(p, &message)
			utils.LogIfError(err2)

			contoller.Service.StoreMessage(message)

			for _, connection := range contoller.Connections {
				if connection.UserId == message.Recipient {
					message, err := json.Marshal(dto.MessageDTO{Sender: message.Sender, Recipient: message.Recipient, Message: message.Message})
					log.Print(err)
					err = connection.Connection.WriteMessage(websocket.TextMessage, message)
					utils.LogIfError(err)
					break
				}
			}
		}
	}(&connection, contoller.ConnectionChan)
}

func (contoller *ChatControllerImpl) GetMessageHistory(ctx *gin.Context) {
	var reqBody dto.GetMessageHistoryDTO
	err := ctx.ShouldBind(&reqBody)
	utils.LogIfError(err)

	history := contoller.Service.GetMessageHistory(reqBody)
	ctx.JSON(200, gin.H{"data": history})
}
