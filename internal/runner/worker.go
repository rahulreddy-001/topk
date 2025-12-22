package runner

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

func StartLeaderboardWorker(
	ctx context.Context,
	rdb *redis.Client,
	consumer *kafka.Reader,
) {
	events := make(chan ScoreEvent, 10_000)

	go func() {
		defer close(events)

		for {
			select {
			case <-ctx.Done():
				log.Println("kafka consumer stopped")
				return

			default:
				msg, err := consumer.ReadMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return
					}
					log.Println("kafka read error:", err)
					continue
				}

				var ev ScoreEvent
				if json.Unmarshal(msg.Value, &ev) == nil {
					events <- ev
				}
			}
		}
	}()

	const (
		maxBatch   = 1000
		flushEvery = 50 * time.Millisecond
	)

	ticker := time.NewTicker(flushEvery)
	defer ticker.Stop()

	batch := make([]ScoreEvent, 0, maxBatch)

	flush := func(batch []ScoreEvent) {
		if len(batch) == 0 {
			return
		}

		pipe := rdb.Pipeline()
		for _, ev := range batch {
			pipe.ZIncrBy(ctx, LB_GLOBAL_REDIS_KEY, float64(ev.Score), ev.UserID)
		}
		_, _ = pipe.Exec(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("leaderboard worker flushing before exit")
			flush(batch)
			batch = batch[:0]
			return

		case ev, ok := <-events:
			if !ok {
				flush(batch)
				batch = batch[:0]
			}
			batch = append(batch, ev)
			if len(batch) >= maxBatch {
				flush(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			flush(batch)
			batch = batch[:0]
		}
	}
}
