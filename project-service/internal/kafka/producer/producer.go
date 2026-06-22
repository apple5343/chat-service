package producer

import (
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

// Producer представляет Kafka Producer
type Producer struct {
	asyncProducer sarama.AsyncProducer
}

// NewProducer создает новый Kafka Producer
func NewProducer(brokers []string, config *sarama.Config) *Producer {
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}

	p := &Producer{asyncProducer: producer}

	// Обработка подтверждений и ошибок
	go func() {
		for {
			select {
			case success := <-p.asyncProducer.Successes():
				log.Info().Msgf("Message sent successfully: topic=%s partition=%d offset=%d",
					success.Topic, success.Partition, success.Offset)
			case err = <-p.asyncProducer.Errors():
				log.Error().Msgf("Failed to send message: %v", err)
			}
		}
	}()

	return p
}

// SendMessage отправляет сообщение в указанный топик
func (p *Producer) SendMessage(topic string, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	p.asyncProducer.Input() <- msg
	return nil
}

// Close закрывает Kafka Producer
func (p *Producer) Close() error {
	return p.asyncProducer.Close()
}
