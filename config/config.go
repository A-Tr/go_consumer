package config

type ServiceConfig struct {
	Port    string `default:":3000"`
	SrvName string `default:"go_consumer"`
}
