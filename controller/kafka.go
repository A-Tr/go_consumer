package controller

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaController struct {
	consumer *kafka.Consumer
}

func InitKafkaController(URL, consumerName string) (*KafkaController, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     URL,
		"broker.address.family": "v4",
		"group.id":              consumerName,
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = c.Subscribe("SOMETOPIC", nil)
	if err != nil {
		return nil, err
	}

	kC := &KafkaController{consumer: c}

	return kC, nil
}

func (c *KafkaController) ConsumeMessages(logger *logrus.Entry) error {
	for {
		ev := c.consumer.Poll(100)

		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			logger.Infof("%% Message on %s:\n%s\n",
				e.TopicPartition, string(e.Value))
			if e.Headers != nil {
				logger.Infof("%% Headers: %v\n", e.Headers)
			}
		case kafka.Error:
			logger.Errorf("Error: %v: %v\n", e.Code(), e)
		default:
			logger.Infof("Ignored %v\n", e)
		}
	}
}
