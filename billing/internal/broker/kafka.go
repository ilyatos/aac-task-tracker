package broker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func Produce(ctx context.Context, topic string, msgs ...kafka.Message) error {
	conn, err := kafka.DialLeader(ctx, "tcp", os.Getenv("KAFKA_URL"), topic, 0)
	if err != nil {
		return fmt.Errorf("failed to dial leader: %w", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	_, err = conn.WriteMessages(msgs...)
	if err != nil {
		return fmt.Errorf("failed to write messages: %w", err)
	}

	if err := conn.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}
