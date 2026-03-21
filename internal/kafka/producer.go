package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/WebCraftersGH/Auth-service/internal/contracts"
	"github.com/WebCraftersGH/Auth-service/internal/domain"
	kafkaGo "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkaGo.Writer
	logger contracts.ILogger
}

func NewProducer(brokers []string, topic string, timeout time.Duration, logger contracts.ILogger) *Producer {
	return &Producer{
		writer: &kafkaGo.Writer{
			Addr:         kafkaGo.TCP(brokers...),
			Topic:        topic,
			RequiredAcks: kafkaGo.RequireOne,
			BatchTimeout: timeout,
		},
		logger: logger,
	}
}

func (p *Producer) PublishUserCreateRequested(ctx context.Context, event domain.UserCreateRequestedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafkaGo.Message{
		Key:   []byte(event.UserID.String()),
		Value: payload,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
