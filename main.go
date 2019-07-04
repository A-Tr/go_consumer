package main

import (
	"go_consumer/config"
	controller "go_consumer/controller"
	log "go_consumer/logger"
	msg "go_consumer/messages"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type App struct {
	BusControllers []controller.BusController
}

var (
	cfg       config.ServiceConfig
	app       App
	loggerCfg log.LoggerConfig
)

func init() {
	log.Init(os.Stdout, logrus.InfoLevel)

	err := envconfig.Process("env", &cfg)
	if err != nil {
		logrus.WithError(err).Fatal(err.Error())
	}

	repository := &msg.MessageMockRepository{}
	mqConsumer, err := controller.InitRabbitController(cfg, repository)
	if err != nil {
		logrus.WithError(err).Fatal("Error connecting to Rabbit Bus")
	}

	kafkaConsumer, err := controller.InitKafkaController("fast-data-dev", "kafka-consumer", repository)
	if err != nil {
		logrus.WithError(err).Fatal("Error connecting to Kafka Bus")
	}

	app = App{
		BusControllers: []controller.BusController{mqConsumer, kafkaConsumer},
	}

}

func main() {

	logger := log.NewLogger(loggerCfg, "CONSUMER_LOGGER")

	logger.Println(" [*] Waiting for logs. To exit press CTRL+C")

	for index := range app.BusControllers {
		go func(bC controller.BusController) {
			for {
				err := bC.ConsumeMessages(logger)
				if err != nil {
					logger.WithError(err).Error("MAIN: Error consuming messages")
				}
			}
		}(app.BusControllers[index])
	}

	logger.Println("Ready to produce and consume")

	// graceful stop
	quitSig := make(chan os.Signal)
	signal.Notify(quitSig, syscall.SIGTERM)
	signal.Notify(quitSig, syscall.SIGINT)
	<-quitSig
	select {
	case sig := <-quitSig:
		logger.Printf("Got %s signal. Handled gracefully.", sig)
		os.Exit(0)
	}
}
