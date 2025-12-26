run:
	sudo docker compose up -d
	go run ./cmd/api/main.go run

flush:
	sudo find ./data -mindepth 2 -delete