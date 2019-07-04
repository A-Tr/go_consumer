package bus

import (
	"bytes"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	msg "go_consumer/messages"
)

type RabbitController struct {
	Q      amqp.Queue
	Ch     *amqp.Channel
	Conn   *amqp.Connection
	Name   string
	busErr error
	mR     msg.MessageRepository
}

func InitRabbitController(srvName string, name string, repository msg.MessageRepository) (*RabbitController, error) {
	config := &RabbitController{
		Name: name,
		mR:   repository,
	}

	config.Conn, config.busErr = amqp.Dial("amqp://guest:guest@172.17.0.2:5672")
	if config.busErr != nil {
		err := errors.Wrapf(config.busErr, "REPO ERROR")
		return nil, err
	}

	config.Ch, config.busErr = config.Conn.Channel()
	if config.busErr != nil {
		err := errors.Wrapf(config.busErr, "REPO ERROR")
		return nil, err
	}

	return config, nil
}

func (b *RabbitController) ConsumeMessages(logger *log.Entry) error {
	for {
		b.Ch, b.busErr = b.Conn.Channel()
		if b.busErr != nil {
			logger.WithError(b.busErr).Error("Error openning channel")
			return b.busErr
		}

		msgs, err := b.Ch.Consume("SOMEQUEUE", b.Name, true, false, false, false, nil)
		if err != nil {
			logger.WithError(err).Error("Error reading messages")
			return err
		}

		for d := range msgs {
			var buffer bytes.Buffer

			buffer.Write(d.Body)
			logger.WithField("consumer: ", b.Name).Infof("Message received %s", d.Body)
			err := b.mR.SaveMessage(d.Body, logger)
			if err != nil {
				logger.WithError(err).Error("Error saving message")
			}
			return err
		}
		defer b.Ch.Close()
	}
}
