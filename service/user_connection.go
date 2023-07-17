package service

import "github.com/gorilla/websocket"

type UserConnection struct {
	UserId     string
	Connection *websocket.Conn
}
