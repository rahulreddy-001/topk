package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func SimulateReads(
	ctx context.Context,
	rdb *redis.Client,
	qps int,
) {
	interval := time.Second / time.Duration(qps)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("read simulator stopped")
			return

		case <-ticker.C:
			b, err := rdb.Get(ctx, LB_GLOBAL_TOP10K_SS_REDIS_KEY).Bytes()
			if err != nil {
				if err != redis.Nil {
					log.Printf("snapshot read error: %v", err)
				}
				continue
			}

			var entries []redis.Z
			if err := json.Unmarshal(b, &entries); err != nil {
				log.Printf("snapshot decode error: %v", err)
				continue
			}

			fmt.Println("----- TOP 10 LEADERBOARD -----")
			limit := 10
			if len(entries) < limit {
				limit = len(entries)
			}

			for i := 0; i < limit; i++ {
				fmt.Printf(
					"#%0.2d user=%v score=%.0f\n",
					i+1,
					entries[i].Member,
					entries[i].Score,
				)
			}
		}
	}
}
