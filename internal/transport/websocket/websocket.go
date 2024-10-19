package websocket

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/yervsil/auth_service/domain"
	mdw "github.com/yervsil/auth_service/internal/transport/http"
)

var ErrUserDisconnected = errors.New("user disconnected")

type Handler struct {
	log 			*slog.Logger
	rooms           map[string]*Room
	Mu         		sync.Mutex
	producer 		mdw.Producer
}

func New(producer mdw.Producer, log *slog.Logger) *Handler {
	return &Handler{producer:producer, 
					log: 	log,
					rooms:  make(map[string]*Room),}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) Chat() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	                                                                                                                 
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			h.log.Error(fmt.Sprintf("error upgrading http to websocket: %s", err.Error()))
			return
		}

		vars := mux.Vars(r)
		roomName := vars["roomId"]

		userData := r.Context().Value(mdw.UserDataKey).(*domain.User)
		
		client := NewClient(h.log, conn, userData, roomName)
		room, err := h.addUserToRoom(client)
		if err != nil {
			h.log.Error(fmt.Sprintf("can't add user to room %s: %v\n", room.Name, err))

			return
		}

		defer h.disconnectUser(client, room)
						
		message := NewMessage(userData.Name, userData.Id, fmt.Sprintf("user %s has joined ",  userData.Name), room.Name) 

		room.PushMessage(message)

		h.handleInputMessages(client, room)
	}
}

func(h *Handler) addUserToRoom(client *Client) (*Room, error) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	room := h.getOrCreateRoom(client.roomName)
	err := room.AddUser(client)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func(h *Handler) getOrCreateRoom(roomName string) *Room {
	if _, ok := h.rooms[roomName]; !ok {
		h.rooms[roomName] = NewRoom(roomName, h.producer)
		go h.rooms[roomName].StartDeliveringMessages()
	}

	room := h.rooms[roomName]

	return room
}

func(h *Handler) disconnectUser(user *Client, room *Room) {
	defer user.conn.Close()

	if room.closed {
		return
	}

	room.DeleteUser(user.user.Id)

	if h.deleteRoomIfEmpty(room) {
		h.log.Info(fmt.Sprintf("room %s was deleted", room.Name))
		return
	}

	msg := NewMessage(user.user.Name, user.user.Id, fmt.Sprintf("user %s left the group", user.user.Name), room.Name) 

	room.PushMessage(msg)
}

func(h *Handler) deleteRoomIfEmpty(room *Room) bool{
	h.Mu.Lock()
	defer h.Mu.Unlock()

	if room.CloseIfEmpty() {
		delete(h.rooms, room.Name)
		return true
	}

	return false
}

func (h *Handler) handleInputMessages(user *Client, room *Room) error {
	for {
		err := h.handleInputMessage(user, room)
		if errors.Is(err, ErrUserDisconnected) {
			return nil
		}

		if err != nil {
			h.log.Error(fmt.Sprintf("can't handle input message: %s", err.Error()))
			return nil
		}
	}
}


func (h *Handler) handleInputMessage(user *Client, room *Room) error {
	_, msg, err := user.conn.ReadMessage()

	if err != nil {
		h.log.Error(fmt.Sprintf("error in reading client WS message: %s", err.Error()))
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			return ErrUserDisconnected
		}

		return err
	}

	message := NewMessage(user.user.Name, user.user.Id, string(msg), room.Name) 

	room.PushMessage(message)

	return nil
}
