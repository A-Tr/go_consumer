package controller

import (
	"bytes"
	"fmt"
	cfg "go_consumer/config"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	msg "go_consumer/messages"
)

type RabbitController struct {
	Q           amqp.Queue
	Ch          *amqp.Channel
	Conn        *amqp.Connection
	Name        string
	busErr      error
	messageRepo msg.MessageRepository
}

func InitRabbitController(config cfg.ServiceConfig, repository msg.MessageRepository) (*RabbitController, error) {
	controller := &RabbitController{
		Name:        config.SrvName,
		messageRepo: repository,
	}

	controller.Conn, controller.busErr = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s",
		config.RabbitUser, config.RabbitPass, config.RabbitURL, config.RabbitPort))
	if controller.busErr != nil {
		err := errors.Wrapf(controller.busErr, "REPO ERROR")
		return nil, err
	}

	controller.Ch, controller.busErr = controller.Conn.Channel()
	if controller.busErr != nil {
		err := errors.Wrapf(controller.busErr, "REPO ERROR")
		return nil, err
	}

	return controller, nil
}

func (b *RabbitController) ConsumeMessages(logger *log.Entry) error {
	for {
		b.Ch, b.busErr = b.Conn.Channel()
		if b.busErr != nil {
			logger.WithError(b.busErr).Error("Error openning channel")
			return b.busErr
		}

		msgs, err := b.Ch.Consume("SOMEQUEUE", b.Name, false, false, false, false, nil)
		if err != nil {
			logger.WithError(err).Error("Error reading messages")
			return err
		}

		for d := range msgs {
			var buffer bytes.Buffer

			buffer.Write(d.Body)
			logger.WithField("consumer: ", b.Name).Infof("Message received %s", d.Body)
			err := b.messageRepo.SaveMessage(d.Body, logger)
			if err != nil {
				logger.WithError(err).Error("Error saving message")
				return err
			}
			d.Ack(true)
		}
		defer b.Ch.Close()
	}
}
