package kafkaHandler

import (
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter returns a new kafka writer
func KafkaWriter(kafkaURL string, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{kafkaURL},
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1,
		BatchTimeout: 10 * time.Millisecond,
	})
}
