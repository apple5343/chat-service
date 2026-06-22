package kafka

import "github.com/IBM/sarama"

type Consumer interface {
	Setup(_ sarama.ConsumerGroupSession) error
	Cleanup(_ sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
	Ready()
	Prepare()
}
