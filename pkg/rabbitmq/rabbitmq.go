package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type ConnectionConfig struct {
	Scheme   string
	User     string
	Password string
	Host     string
	Port     string
	Vhost    string
}

func (r *ConnectionConfig) String() string {

	return fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		r.Scheme,
		r.User,
		r.Password,
		r.Host,
		r.Port,
		r.Vhost,
	)
}

// ConnAndCreateChan пытается подключиться к RabbitMQ и создать канал, с несколькими попытками
func ConnAndCreateChan(cfg ConnectionConfig, maxRetries int, timeout time.Duration) (*amqp.Channel, error) {
	const op = "rabbitmq.ConnAndCreateChan"
	var conn *amqp.Connection
	var err error

	for retries := 0; retries < maxRetries; retries++ {
		conn, err = amqp.Dial(cfg.String())
		if err == nil {
			break
		}
		if retries < maxRetries-1 {
			fmt.Printf("%s: failed to connect to RabbitMQ, retrying in %v... (attempt %d/%d)\n", op, timeout, retries+1, maxRetries)
			time.Sleep(timeout)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to RabbitMQ after %d attempts: %w", op, maxRetries, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create channel: %w", op, err)
	}

	return ch, nil
}

// RabbitMQWriter будет отправлять сообщения в rabbitMQ
type RabbitMQWriter struct {
	Channel    *amqp.Channel
	Exchange   string           `yaml:"exchange"`
	RoutingKey string           `yaml:"routing_key"`
	Mandatory  bool             `yaml:"mandatory"`
	Immediate  bool             `yaml:"immediate"`
	Publishing PublishingConfig `yaml:"publishing"`
}

type PublishingConfig struct {
	ContentType     string `yaml:"content_type"`
	ContentEncoding string `yaml:"content_encoding"`
	DeliveryMode    uint8  `yaml:"delivery_mode"`
	Priority        uint8  `yaml:"priority"`
	CorrelationId   string `yaml:"correlation_id"`
	ReplyTo         string `yaml:"reply_to"`
	Expiration      string `yaml:"expiration"`
	Type            string `yaml:"type"`
	UserId          string `yaml:"user_id"`
	AppId           string `yaml:"app_id"`
}

// Write реализует интерфейс io.Writer
func (r RabbitMQWriter) Write(p []byte) (n int, err error) {
	const op = "rabbitmq.Write"
	// инфо для os.stdout
	log.Println(string(p))

	msg := amqp.Publishing{
		ContentType:     r.Publishing.ContentType,
		ContentEncoding: r.Publishing.ContentEncoding,
		DeliveryMode:    r.Publishing.DeliveryMode,
		Priority:        r.Publishing.Priority,
		CorrelationId:   r.Publishing.CorrelationId,
		ReplyTo:         r.Publishing.ReplyTo,
		Expiration:      r.Publishing.Expiration,
		Type:            r.Publishing.Type,
		UserId:          r.Publishing.UserId,
		AppId:           r.Publishing.AppId,
		Timestamp:       time.Now(),
		Body:            p, // Основное сообщение
	}

	err = r.Channel.Publish(
		r.Exchange,
		r.RoutingKey,
		r.Mandatory,
		r.Immediate,
		msg,
	)
	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("%s: failed to write: %w", op, err)
	}
	return len(p), nil
}
