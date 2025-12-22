package cmd

import (
	"log"
	"topk/internal/runner"

	"github.com/spf13/cobra"
)

func startRunner() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "start runner",
		Run: func(_ *cobra.Command, _ []string) {
			// init config
			cfg, err := LoadConfig()
			if err != nil {
				log.Fatalf("error loading config: %v", err)
			}
			log.Println("loaded config successfully")

			// init redis
			rdb, err := getRedisClient(cfg)
			if err != nil {
				log.Fatalf("error initilizing redis client: %v", err)
			}
			log.Println("initilized redis client successfully")

			// init kafka
			kafkaWriter, kafkaReader, err := initKafka(cfg)
			if err != nil {
				log.Fatalf("error initilizing kafka: %v", err)
			}
			log.Println("initilized kafka successfully")

			// start runner
			runner.Start(rdb, kafkaWriter, kafkaReader)

		},
	}
}
