package config

import (
	"os"
	"strings"

	"github.com/IBM/sarama"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	topicsEnvName  = "KAFKA_TOPICS"
	groupIDEnvName = "KAFKA_GROUP_ID"
)

type KafkaConfig interface {
	SaramaConfig() *sarama.Config
	Brokers() []string
	Topics() []string
	GroupID() string
}

type kafkaConfig struct {
	saramaConfig *sarama.Config
	brokers      []string
	topics       []string
	groupID      string
}

func NewKafkaConfig() KafkaConfig {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Version = sarama.V3_5_0_0

	brokers := strings.Split(os.Getenv(brokersEnvName), " ")
	topics := strings.Split(os.Getenv(topicsEnvName), " ")
	groupID := os.Getenv(groupIDEnvName)

	return &kafkaConfig{
		saramaConfig: config,
		brokers:      brokers,
		topics:       topics,
		groupID:      groupID,
	}
}

func (cfg *kafkaConfig) SaramaConfig() *sarama.Config {
	return cfg.saramaConfig
}

func (cfg *kafkaConfig) Brokers() []string {
	return cfg.brokers
}

func (cfg *kafkaConfig) Topics() []string {
	return cfg.topics
}

func (cfg *kafkaConfig) GroupID() string {
	return cfg.groupID
}
