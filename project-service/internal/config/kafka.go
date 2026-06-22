package config

import (
	"os"
	"strings"

	"github.com/IBM/sarama"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	topicEnvName   = "KAFKA_TOPIC"
)

type KafkaConfig interface {
	SaramaConfig() *sarama.Config
	Brokers() []string
	Topic() string
}

type kafkaConfig struct {
	saramaConfig *sarama.Config
	brokers      []string
	topic        string
}

func NewKafkaConfig() KafkaConfig {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Version = sarama.V3_5_0_0

	brokers := strings.Split(os.Getenv(brokersEnvName), " ")
	topic := os.Getenv(topicEnvName)

	return &kafkaConfig{
		saramaConfig: config,
		brokers:      brokers,
		topic:        topic,
	}
}

func (cfg *kafkaConfig) SaramaConfig() *sarama.Config {
	return cfg.saramaConfig
}

func (cfg *kafkaConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaConfig) Topic() string {
	return cfg.topic
}
