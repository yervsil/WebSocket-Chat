package websocket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"log/slog"
	"github.com/yervsil/auth_service/domain"
	mdw "github.com/yervsil/auth_service/internal/transport/http"
)

type Room struct {
	Name    		string
	Mu 		        sync.RWMutex
	closed  		bool 
	users 			map[int]*Client
	messageStack 	chan Message
	producer 		mdw.Producer
}

func NewRoom(roomName string, producer mdw.Producer) *Room {
	return &Room{
				Name: roomName,
				messageStack: make(chan Message),
				users: make(map[int]*Client),
				producer: producer,
			}
}

func (r *Room) AddUser(client *Client) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if r.closed {
		return fmt.Errorf("romm %s was closed", r.Name)
	}

	r.users[client.user.Id] = client

	return nil
}

func (r *Room) DeleteUser(userID int) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	delete(r.users, userID)
}

func (r *Room) PushMessage(msg Message) {
	r.messageStack <- msg
}

func (r *Room) StartDeliveringMessages() {
	for msg := range r.messageStack {
		r.sendMessageToOthers(msg)
	}
}

func (r *Room) CloseIfEmpty() bool {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if len(r.users) == 0 {
		close(r.messageStack)
		r.closed = true

		return true
	}

	return false
}

func (r *Room) sendMessageToOthers(msg Message) {
	r.Mu.RLock()
	defer r.Mu.RUnlock()

	kafkaMessages := make([]domain.KafkaLoggingMessage, 0, len(r.users))

	for _, u := range r.users {
		if u.user.Id != msg.UserId {
			err := u.conn.WriteMessage(1, []byte(fmt.Sprintf("%s: %s", msg.UserName, msg.Message)))
			if err != nil {
				slog.Error(fmt.Sprintf("can't send messages to kafka: %s", err.Error()))
				
			} else {
				kafkaMessages = append(kafkaMessages, newKafkaLoggingMessage(msg, u.user.Id))
			}
	}
		}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := r.producer.ProduceJSONMessage(ctx, kafkaMessages)
		if err != nil {
			slog.Error(fmt.Sprintf("can't send messages to kafka: %s", err.Error()))
		}
		slog.Info(fmt.Sprintf("message %s has been produced", msg.Message))
	}()
}

func newKafkaLoggingMessage(msg Message, toUserId int) domain.KafkaLoggingMessage {
	messageUUID := uuid.New()

	return domain.KafkaLoggingMessage{
		MessageUUID:    messageUUID,
		Timestamp:      time.Now(),
		MessageContent: string(msg.Message),
		FromUserUUID:   msg.UserId,
		ToUserUUID:     toUserId,
	}
}