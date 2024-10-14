package websocket

import (
	"log/slog"

	"github.com/gorilla/websocket"
	"github.com/yervsil/auth_service/domain"
)

type Client struct {
	log      *slog.Logger
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

func NewClient(log *slog.Logger, conn *websocket.Conn,
	user *domain.User, roomName  string) *Client {
	return &Client{log: log, conn: conn, user: user, roomName: roomName}
}
