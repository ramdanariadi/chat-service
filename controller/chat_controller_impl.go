package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ramdanariadi/chat-service/dto"
	"github.com/ramdanariadi/chat-service/service"
	"github.com/ramdanariadi/chat-service/utils"
	"log"
)

type ChatControllerImpl struct {
	Upgrader       *websocket.Upgrader
	Connections    *[]*service.UserConnection
	ConnectionChan chan *service.UserConnection
	Service        service.ChatService
}

func (controller *ChatControllerImpl) WSHandler(ctx *gin.Context) {
	conn, err := controller.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	utils.LogIfError(err)
	userId := ctx.Param("userId")
	connection := service.UserConnection{UserId: userId, Connection: conn}
	*controller.Connections = append(*controller.Connections, &connection)
	marshal, err := json.Marshal(dto.MessageDTO{Recipient: userId, Message: "hi!!"})
	utils.LogIfError(err)

	err = conn.WriteMessage(websocket.TextMessage, marshal)
	log.Println(err)

	go func(userConnection *service.UserConnection, userChan chan<- *service.UserConnection) {
		for {
			_, p, e := userConnection.Connection.ReadMessage()
			if nil != e {
				log.Println("read message from " + userConnection.UserId)
				utils.LogIfError(e)
				if websocket.IsUnexpectedCloseError(e) {
					log.Println("ws closed for " + userConnection.UserId)
					userChan <- userConnection
				}
				break
			}
			log.Println(p)
			var message dto.MessageDTO
			err2 := json.Unmarshal(p, &message)
			utils.LogIfError(err2)

			controller.Service.StoreMessage(message)
			log.Printf("connection base in WS HANDLER : %p", controller.Connections)
			for _, connection := range *controller.Connections {
				log.Printf("check connection IN WS HANDLER,userId : %s,connection addr : %p,userId addr : %p, ws.conn Addr : %p", connection.UserId, connection, &connection.UserId, connection.Connection)
				if connection.UserId == message.Recipient {
					message, err := json.Marshal(dto.MessageDTO{Sender: message.Sender, Recipient: message.Recipient, Message: message.Message})
					log.Print(err)
					err = connection.Connection.WriteMessage(websocket.TextMessage, message)
					utils.LogIfError(err)
					break
				}
			}
		}
	}(&connection, controller.ConnectionChan)
}

func (controller *ChatControllerImpl) GetMessageHistory(ctx *gin.Context) {
	var reqBody dto.GetMessageHistoryDTO
	err := ctx.ShouldBind(&reqBody)
	utils.LogIfError(err)

	history := controller.Service.GetMessageHistory(reqBody)
	ctx.JSON(200, gin.H{"data": history})
}

func (controller *ChatControllerImpl) GetUser(ctx *gin.Context) {
	var requestBody dto.UserChatRequest
	err := ctx.ShouldBind(&requestBody)
	utils.LogIfError(err)
	message := controller.Service.GetUserWithLastMessage(requestBody)
	ctx.JSON(200, gin.H{"data": message})
}
