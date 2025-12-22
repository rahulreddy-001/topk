package cmd

import (
	"topk/internal/utils"
)

type Config struct {
	Database struct {
		Redis struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			Password string `json:"password"`
			DB       int    `json:"db"`
		} `json:"redis"`
	} `json:"database"`

	Kafka struct {
		Brokers []string `json:"brokers"`
		Topic   string   `json:"topic"`
		GroupID string   `json:"group_id"`
	} `json:"kafka"`
}

func LoadConfig() (*Config, error) {
	return utils.LoadJSONFromFile[Config]("config.json")
}
