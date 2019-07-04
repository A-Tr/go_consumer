package messages

import (
	"time"
)

type Message struct {
	Message   string    `json:"Message"`
	Publisher string    `json:"Publisher"`
	Timestamp time.Time `json:"Timestamp"`
}
