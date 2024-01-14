migrate:
	go run ./cmd/migrator --storage-path=./storage/cards.db --migrations-path=./migrations

run:
	go run ./cmd/cards --config=./config/config.yaml
