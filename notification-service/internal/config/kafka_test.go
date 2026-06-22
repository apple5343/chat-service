package config_test

import (
	"github.com/IBM/sarama"
	"notification-service/internal/config"
	"os"
	"reflect"
	"testing"
)

func TestNewKafkaConfig(t *testing.T) {
	// Устанавливаем переменные окружения
	os.Setenv("KAFKA_BROKERS", "broker1 broker2")
	os.Setenv("KAFKA_TOPICS", "topic1 topic2")
	os.Setenv("KAFKA_GROUP_ID", "group1")
	defer func() {
		os.Unsetenv("KAFKA_BROKERS")
		os.Unsetenv("KAFKA_TOPICS")
		os.Unsetenv("KAFKA_GROUP_ID")
	}()

	// Создаем конфигурацию
	kafkaCfg := config.NewKafkaConfig()

	// Проверяем SaramaConfig
	saramaCfg := kafkaCfg.SaramaConfig()
	if saramaCfg.Producer.RequiredAcks != sarama.WaitForAll {
		t.Errorf("Expected RequiredAcks to be %v, got %v", sarama.WaitForAll, saramaCfg.Producer.RequiredAcks)
	}
	if saramaCfg.Producer.Retry.Max != 5 {
		t.Errorf("Expected Retry.Max to be 5, got %v", saramaCfg.Producer.Retry.Max)
	}
	if !saramaCfg.Producer.Return.Successes {
		t.Errorf("Expected Return.Successes to be true, got %v", saramaCfg.Producer.Return.Successes)
	}
	if saramaCfg.Version != sarama.V3_5_0_0 {
		t.Errorf("Expected Sarama version to be %v, got %v", sarama.V3_5_0_0, saramaCfg.Version)
	}

	expectedBrokers := []string{"broker1", "broker2"}
	if !reflect.DeepEqual(kafkaCfg.Brokers(), expectedBrokers) {
		t.Errorf("Expected brokers to be %v, got %v", expectedBrokers, kafkaCfg.Brokers())
	}

	expectedTopics := []string{"topic1", "topic2"}
	if !reflect.DeepEqual(kafkaCfg.Topics(), expectedTopics) {
		t.Errorf("Expected topics to be %v, got %v", expectedTopics, kafkaCfg.Topics())
	}

	expectedGroupID := "group1"
	if kafkaCfg.GroupID() != expectedGroupID {
		t.Errorf("Expected group ID to be %v, got %v", expectedGroupID, kafkaCfg.GroupID())
	}
}

func TestKafkaConfigWithoutEnv(t *testing.T) {
	os.Unsetenv("KAFKA_BROKERS")
	os.Unsetenv("KAFKA_TOPICS")
	os.Unsetenv("KAFKA_GROUP_ID")

	kafkaCfg := config.NewKafkaConfig()

	if len(kafkaCfg.Brokers()) != 1 || kafkaCfg.Brokers()[0] != "" {
		t.Errorf("Expected brokers to be empty, got %v", kafkaCfg.Brokers())
	}
	if len(kafkaCfg.Topics()) != 1 || kafkaCfg.Topics()[0] != "" {
		t.Errorf("Expected topics to be empty, got %v", kafkaCfg.Topics())
	}
	if kafkaCfg.GroupID() != "" {
		t.Errorf("Expected group ID to be empty, got %v", kafkaCfg.GroupID())
	}
}
