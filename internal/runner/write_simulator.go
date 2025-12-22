package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

func SimulateScores(
	ctx context.Context,
	writer *kafka.Writer,
	userCount int,
	eventsPerSec int,
) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("score simulator stopped")
			return

		case <-ticker.C:
			msgs := make([]kafka.Message, 0, eventsPerSec)

			for i := 0; i < eventsPerSec; i++ {
				ev := ScoreEvent{
					UserID: fmt.Sprintf("user_%d", rand.Intn(userCount)),
					Score:  int64(rand.Intn(50) + 1),
					Ts:     time.Now().UnixMilli(),
				}

				b, _ := json.Marshal(ev)

				msgs = append(msgs, kafka.Message{
					Key:   []byte(ev.UserID),
					Value: b,
				})
			}

			if err := writer.WriteMessages(ctx, msgs...); err != nil {
				log.Println("produce error:", err)
			}
		}
	}
}
