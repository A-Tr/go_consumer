package messages

import (
	"github.com/sirupsen/logrus"
	"encoding/json"
)

type MessageMockRepository struct{}

func (m *MessageMockRepository) SaveMessage(msgBytes []byte, logger *logrus.Entry) error {
	msg := Message{}

	err := json.Unmarshal(msgBytes, &msg)
	if err != nil {
		logger.WithError(err).Errorf("Error unmarshalling msg. Raw message: %s", msgBytes)
		return err
	}

	logger.Infof("Message converted to struct %+v", msg)
	return nil
}
