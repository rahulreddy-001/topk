package runner

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type ScoreEvent struct {
	UserID string `json:"user_id"`
	Score  int64  `json:"score"`
	Ts     int64  `json:"ts"`
}

const (
	LB_GLOBAL_REDIS_KEY           = "LB_GLOBAL"
	LB_GLOBAL_TOP10K_SS_REDIS_KEY = "LB_GLOBAL_TOP10K_SS"
)

func Start(
	rdb *redis.Client,
	producer *kafka.Writer,
	consumer *kafka.Reader,
) {
	log.Println("starting runner...")

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Go(func() { SimulateScores(ctx, producer, 1_000_000, 1_000_000) })
	wg.Go(func() { StartLeaderboardWorker(ctx, rdb, consumer) })
	wg.Go(func() { StartSnapshotBuilder(ctx, rdb) })
	wg.Go(func() { SimulateReads(ctx, rdb, 20_000) })

	shutdown()
	log.Println("shutdown signal received")
	cancel()

	_ = producer.Close()
	_ = consumer.Close()

	wg.Wait()

	log.Println("runner shutdown complete")
}

func shutdown() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sig)
	<-sig
}
