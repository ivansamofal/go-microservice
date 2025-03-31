.PHONY: db-logs db go build deps

db-logs:
	docker logs -f go-postgres

db:
	docker exec -it go-postgres bash

go:
	docker exec -it go-microservice bash

go-build:
	go build -o main .

deps:
	go mod tidy

build:
	docker compose up -d --build && docker exec node-frontend npm i && docker exec node-frontend npm run build

npm:
	docker exec -it node-frontend bash