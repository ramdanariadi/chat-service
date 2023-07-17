package main

import (
	"github.com/gin-contrib/cors"
	_ "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/ramdanariadi/chat-service/controller"
	"github.com/ramdanariadi/chat-service/model"
	"github.com/ramdanariadi/chat-service/service"
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

var connections = make([]*service.UserConnection, 0)
var connChan = make(chan *service.UserConnection)

func CloseUserConnection(connChan <-chan *service.UserConnection) {
	var connTemp = make([]*service.UserConnection, 0)
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

	chatController := controller.ChatControllerImpl{DB: db, Upgrader: &upgrader, ConnectionChan: connChan, Connections: connections}
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

	router.GET("/ws/:userId", chatController.WSHandler)
	router.POST("/message/history", chatController.GetMessageHistory)
	router.POST("/message/history/user")

	err = router.Run(":8087")
	utils.LogIfError(err)
}
