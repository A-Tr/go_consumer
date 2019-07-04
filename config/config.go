package config

type RabbitConfig struct {
	RabbitUser string `default:"guest"`
	RabbitPass string `default:"guest"`
	RabbitURL  string `default:"localhost"`
	RabbitPort string `default:"5672"`
}
type ServiceConfig struct {
	Port    string `default:":3000"`
	SrvName string `default:"go_consumer"`
	RabbitConfig
}
