package main

import (
	"encoding/json"
	"github.com/gin-contrib/cors"
	_ "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/ramdanariadi/chat-service/dto"
	"github.com/ramdanariadi/chat-service/model"
	"github.com/ramdanariadi/chat-service/setup"
	"github.com/ramdanariadi/chat-service/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type UserConnection struct {
	UserId     string
	Connection *websocket.Conn
}

var connections = make([]*UserConnection, 0)
var connChan = make(chan *UserConnection)

func (contoller *ChatControllerImpl) wsHandlerNew(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	userId := ctx.Param("userId")
	connection := UserConnection{UserId: userId, Connection: conn}
	connections = append(connections, &connection)
	marshal, err := json.Marshal(dto.MessageDTO{Recipient: userId, Message: "hi!!"})

	if nil != err {
		log.Println(err)
	}

	err = conn.WriteMessage(websocket.TextMessage, marshal)
	log.Println(err)

	go func(userConnection *UserConnection, userChan chan<- *UserConnection) {
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
			log.Println("message : " + message.Message)

			newUUID, _ := uuid.NewUUID()
			chat := model.Chat{Id: newUUID.String(), Sender: message.Sender, Recipient: message.Recipient, Message: message.Message}
			contoller.DB.Save(&chat)

			for _, connection := range connections {
				if connection.UserId == message.Recipient {
					message, err := json.Marshal(dto.MessageDTO{Recipient: userConnection.UserId, Message: message.Message})
					log.Print(err)
					err = connection.Connection.WriteMessage(websocket.TextMessage, message)
					utils.LogIfError(err)
					break
				}
			}
			log.Println("end")
		}
	}(&connection, connChan)
}

func CloseUserConnection(connChan <-chan *UserConnection) {
	var connTemp = make([]*UserConnection, 0)
	for userConn := range connChan {
		log.Println("CloseUserConnection")
		for _, conn := range connections {
			if conn == userConn {
				err := conn.Connection.Close()
				utils.LogIfError(err)
			} else {
				connTemp = append(connTemp, conn)
			}
		}
		connections = connTemp
	}
}

type ChatController interface {
	getMessageHistory(ctx *gin.Context)
	wsHandlerNew(ctx *gin.Context)
}

type ChatControllerImpl struct {
	gorm.DB
}

func (contoller *ChatControllerImpl) getMessageHistory(ctx *gin.Context) {
	var reqBody dto.GetMessageHistoryDTO
	chatHistory := make([]dto.MessageHistoryDTO, 0)
	var chats []*model.Chat

	err := ctx.ShouldBind(&reqBody)
	utils.LogIfError(err)
	senderAndRecipient := []string{reqBody.UserIdFrom, reqBody.UserIdTo}
	tx := contoller.DB.Model(&model.Chat{}).Where("sender IN ? AND recipient IN ?", senderAndRecipient, senderAndRecipient).Find(&chats)
	utils.LogIfError(tx.Error)

	for _, chat := range chats {
		dto := dto.MessageHistoryDTO{From: chat.Sender, To: chat.Recipient, Message: chat.Message, Time: chat.CreatedAt.Unix()}
		chatHistory = append(chatHistory, dto)
	}

	ctx.JSON(200, gin.H{"data": chatHistory})
}

func main() {
	connection, err := setup.NewDBConnection()
	utils.LogIfError(err)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: connection}))
	utils.LogIfError(err)
	err = db.AutoMigrate(model.Chat{})
	utils.LogIfError(err)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	go CloseUserConnection(connChan)
	defer func() {
		close(connChan)
	}()

	chatController := ChatControllerImpl{DB: *db}
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/ws/:userId", chatController.wsHandlerNew)
	router.POST("/message", chatController.getMessageHistory)

	err = router.Run(":8087")
	utils.LogIfError(err)
}
