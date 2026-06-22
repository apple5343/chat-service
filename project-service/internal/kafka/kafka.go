package kafka

type Producer interface {
	SendMessage(topic string, message string) error
	Close() error
}
