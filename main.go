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
	BusController controller.BusController
}

var (
	cfg       config.ServiceConfig
	app       App
	loggerCfg log.LoggerConfig
)

func init() {
	log.Init(os.Stdout, logrus.InfoLevel)

	err := envconfig.Process("local", &cfg)
	if err != nil {
		logrus.WithError(err).Fatal(err.Error())
	}

	repository := &msg.MessageMockRepository{}
	mqConsumer, err := controller.InitRabbitController(cfg.SrvName, "consumer", repository)
	if err != nil {
		logrus.WithError(err).Fatal("Error connecting to bus")
	}

	app = App{
		BusController: mqConsumer,
	}

}

func main() {

	logger := log.NewLogger(loggerCfg, "CONSUMER_LOGGER")
	logger.Print("Starting server " + cfg.SrvName + " on port " + cfg.Port)

	logger.Println(" [*] Waiting for logs. To exit press CTRL+C")

	go func() {
		for {
			err := app.BusController.ConsumeMessages(logger)
			if err != nil {
				logger.WithError(err).Error("MAIN: Error consuming messages")
			}
		}

	}()

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
