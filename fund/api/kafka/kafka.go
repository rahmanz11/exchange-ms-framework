package kafkaHandler

import (
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter creates a new kafka writer
func KafkaWriter(kafkaURL string, topic string) *kafka.Writer {
	// return &kafka.Writer{
	// 	Addr:     kafka.TCP(kafkaURL),
	// 	Topic:    topic,
	// 	Balancer: &kafka.LeastBytes{},
	// }

	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{kafkaURL},
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1,
		BatchTimeout: 10 * time.Millisecond,
	})
}

// KafkaReader creates a new kafka reader
func KafkaReader(kafkaURL, groupID, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,    // 10KB
		MaxBytes: 10e6, // 10MB
	})
}
