package domain

import (
	"time"

	"github.com/google/uuid"
)

type KafkaLoggingMessage struct {
	MessageUUID    uuid.UUID
	Timestamp      time.Time
	MessageContent string
	FromUserUUID   int
	ToUserUUID     int
}