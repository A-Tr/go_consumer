package config

type RabbitConfig struct {
	RabbitUser string `default:"guest"`
	RabbitPass string `default:"guest"`
	RabbitURL  string `default:"localhost"`
	RabbitPort string `default:"5672"`
}

type KafkaConfig struct {
	KafkaUrl      string `default:"fast-data-dev"`
	KafkaConsumer string `default:"kafka-consumer"`
	KafkaTopic    string `default:"SOMETOPIC"`
}

type ServiceConfig struct {
	SrvName string `default:"go_consumer"`
	RabbitConfig
	KafkaConfig
}
