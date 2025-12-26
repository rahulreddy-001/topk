package cmd

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

func initKafka(cfg *Config) (*kafka.Writer, *kafka.Reader, error) {

	kafkaWriter := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Brokers...),
		Topic:        cfg.Kafka.Topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireOne,
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		Async:        false,
	}

	kafkaReader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        cfg.Kafka.Brokers,
			Topic:          cfg.Kafka.Topic,
			GroupID:        cfg.Kafka.GroupID,
			MinBytes:       1024,
			MaxBytes:       10 * 1024 * 1024,
			MaxWait:        500 * time.Millisecond,
			CommitInterval: time.Second,
			StartOffset:    kafka.FirstOffset,
		},
	)
	connCheck := kafkaPing(cfg.Kafka.Brokers)
	if connCheck != nil {
		return  nil, nil, connCheck
	}
	err := ensureTopic(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		return  nil, nil, err
	}
	return kafkaWriter, kafkaReader, connCheck
}

func kafkaPing(brokers []string) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Brokers()
	return err
}

func ensureTopic(brokers []string, topic string) error {
	client := kafka.Client{
		Addr: kafka.TCP(brokers...),
	}

	req := &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     3,
				ReplicationFactor: 3,
			},
		},
	}

	_, err := client.CreateTopics(context.Background(), req)
	return err
}
