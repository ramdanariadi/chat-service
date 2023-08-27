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
	//marshal, err := json.Marshal(dto.MessageDTO{Recipient: userId, Message: "hi!!"})
	//utils.LogIfError(err)

	//err = conn.WriteMessage(websocket.TextMessage, marshal)
	//log.Println(err)

	go func(userConnection *service.UserConnection, userChan chan<- *service.UserConnection) {
		for {
			_, p, e := userConnection.Connection.ReadMessage()
			utils.LogIfError(e)
			log.Println("read messageDTO from userId : " + userConnection.UserId)
			if nil != e {
				if websocket.IsUnexpectedCloseError(e) {
					log.Println("ws closed for userId : " + userConnection.UserId)
					userChan <- userConnection
				}
				break
			}
			log.Println(p)
			var messageDTO dto.MessageDTO
			err2 := json.Unmarshal(p, &messageDTO)
			utils.LogIfError(err2)

			controller.Service.StoreMessage(messageDTO)
			log.Printf("connection base in WS Handler => %p, len : %d", controller.Connections, len(*controller.Connections))
			for index, connection := range *controller.Connections {
				log.Printf("%d. Connection Info => userId : %s, connection addr : %p, userId addr : %p, ws.conn addr : %p", index, connection.UserId, connection, &connection.UserId, connection.Connection)
				if connection.UserId == messageDTO.Recipient {
					message, err := json.Marshal(messageDTO)
					utils.LogIfError(err)
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
	ctx.JSON(200, history)
}

func (controller *ChatControllerImpl) GetUser(ctx *gin.Context) {
	var requestBody dto.UserChatRequest
	err := ctx.ShouldBind(&requestBody)
	utils.LogIfError(err)
	message := controller.Service.GetUserWithLastMessage(requestBody)
	ctx.JSON(200, gin.H{"data": message})
}

func (controller *ChatControllerImpl) CloseUserConnection() {
	for userConn := range controller.ConnectionChan {
		var connTemp = make([]*service.UserConnection, 0)
		for _, conn := range *controller.Connections {
			if conn == userConn {
				err := conn.Connection.Close()
				utils.LogIfError(err)
				log.Printf("CloseUserConnection => userId : %s,connection addr : %p,userId addr : %p, ws.conn Addr : %p", conn.UserId, conn, &conn.UserId, conn.Connection)
			} else {
				connTemp = append(connTemp, conn)
			}
		}
		*controller.Connections = connTemp
		log.Printf("User Connections after remove closed connection : %d", len(*controller.Connections))
		for _, conn := range *controller.Connections {
			log.Printf("User Connection : %s, address : %p", conn.UserId, conn)
		}
	}
}
