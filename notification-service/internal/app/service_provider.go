package app

import (
	"fmt"
	"notification-service/internal/closer"
	"notification-service/internal/config"
	"notification-service/internal/kafka"
	"notification-service/internal/kafka/consumer"
	"notification-service/internal/telegram"
	"notification-service/internal/telegram/notifier"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
)

type serviceProvider struct {
	loggerConfig   config.LoggerConfig
	kafkaConfig    config.KafkaConfig
	notifierConfig config.NotifierConfig

	consumer            kafka.Consumer
	consumerGroup       sarama.ConsumerGroup
	notificationHandler consumer.NotificationHandler

	notifier telegram.Notifier
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) LoggerConfig() config.LoggerConfig {
	if s.loggerConfig == nil {
		cfg, err := config.NewLoggerConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get logger config: %w", err))
		}

		s.loggerConfig = cfg
	}

	return s.loggerConfig
}

func (s *serviceProvider) KafkaConfig() config.KafkaConfig {
	if s.kafkaConfig == nil {
		cfg := config.NewKafkaConfig()
		s.kafkaConfig = cfg
	}

	return s.kafkaConfig
}

func (s *serviceProvider) NotifierConfig() config.NotifierConfig {
	if s.notifierConfig == nil {
		cfg, err := config.NewNotifierConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get notifier config: %w", err))
		}

		s.notifierConfig = cfg
	}

	return s.notifierConfig
}

func (s *serviceProvider) NotificationHandler() consumer.NotificationHandler {
	if s.notificationHandler == nil {
		s.notificationHandler = func(chatID int64, message string) {
			log.Info().Msg(message)
			err := s.Notifier().Notify(chatID, message)
			if err != nil {
				log.Error().Err(err).Msg("failed to send notification")
			}
		}
	}

	return s.notificationHandler
}

func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = consumer.NewConsumer(s.NotificationHandler())
	}

	return s.consumer
}

func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(s.KafkaConfig().Brokers(), s.KafkaConfig().GroupID(), s.KafkaConfig().SaramaConfig())
		if err != nil {
			panic(fmt.Errorf("failed to create consumer group: %w", err))
		}
		closer.Add(consumerGroup.Close)

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

func (s *serviceProvider) Notifier() telegram.Notifier {
	if s.notifier == nil {
		n, err := notifier.NewNotifier(s.NotifierConfig().BotToken())
		if err != nil {
			panic(fmt.Errorf("failed to get notifier: %w", err))
		}

		s.notifier = n
	}

	return s.notifier
}
