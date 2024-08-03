package msgbroker

import (
	"NoteProject/pkg/rabbitmq"
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/streadway/amqp"
)

func InitWriterMustLoad(ch *amqp.Channel) rabbitmq.RabbitMQWriter {
	configPath := filepath.Join("config", "rabbitMQ.yaml")

	if configPath == "" {
		log.Fatal("invalid path to config")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("cfg file not found: %s", configPath)
	}

	cfg := rabbitmq.RabbitMQWriter{Channel: ch}
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cfg file is incorrect: %s", err)
	}
	return cfg

}
