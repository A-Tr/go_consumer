package messages

import (
	"github.com/sirupsen/logrus"
)

type MessageRepository interface {
	SaveMessage([]byte, *logrus.Entry) error
}
