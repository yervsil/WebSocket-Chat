package websocket

import (
	"log/slog"

	"github.com/gorilla/websocket"
	"github.com/yervsil/auth_service/domain"
)

type Client struct {
	conn     *websocket.Conn
	user 	 *domain.User
	roomName  string 
}

type Message struct {
	UserName    string 	`json:"username"`
	UserId	    int		`json:"userId"`
	Message	    string  `json:"message"`
	RoomName    string  `json:"roomName"`
}

func NewMessage(username string, userid int, message string, roomName string) Message {
	return Message{UserName: username, UserId: userid, Message: message, RoomName: roomName}
}

func NewClient(log *slog.Logger, conn *websocket.Conn,
	user *domain.User, roomName  string) *Client {
	return &Client{conn: conn, user: user, roomName: roomName}
}
