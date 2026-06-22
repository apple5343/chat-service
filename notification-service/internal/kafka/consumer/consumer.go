package consumer

import (
	"strconv"
	"strings"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

// NotificationHandler обрабатывает полученные сообщения
type NotificationHandler func(chatID int64, message string)

// Consumer представляет Kafka Consumer
type Consumer struct {
	ready  chan bool
	handle NotificationHandler
}

// NewConsumer создает новый Consumer
func NewConsumer(handler NotificationHandler) *Consumer {
	return &Consumer{
		ready:  make(chan bool),
		handle: handler,
	}
}

// Setup запускается перед началом обработки сообщений
func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup вызывается после завершения обработки
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim обрабатывает сообщения из топика
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Info().Msgf("Message received: topic=%s partition=%d offset=%d value=%s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		parts := strings.Split(string(msg.Value), ":-:")
		chatID, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing chat ID")
		}

		c.handle(chatID, parts[1]) // Передаем сообщение в обработчик
		session.MarkMessage(msg, "")
	}
	return nil
}

func (c *Consumer) Prepare() {
	c.ready = make(chan bool)
}

func (c *Consumer) Ready() {
	<-c.ready
}
