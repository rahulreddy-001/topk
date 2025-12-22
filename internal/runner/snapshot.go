package runner

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func StartSnapshotBuilder(
	ctx context.Context,
	rdb *redis.Client,
) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("snapshot builder stopped")
			return

		case <-ticker.C:
			res, err := rdb.ZRevRangeWithScores(
				ctx,
				LB_GLOBAL_REDIS_KEY,
				0,
				9999,
			).Result()

			if err != nil {
				continue
			}

			b, _ := json.Marshal(res)
			rdb.Set(ctx, LB_GLOBAL_TOP10K_SS_REDIS_KEY, b, time.Second)
		}
	}
}
