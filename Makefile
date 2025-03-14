QUIET := @
ARGS=$(filter-out $@, $(MAKECMDGOALS))

.DEFAULT_GOAL=help
.PHONY=help
app_container=app
queue_container=app
app_container=queue
db_container=postgresql

db-logs:
    docker logs -f go-postgres

db:
    docker exec -it go-postgres bash

go:
    docker exec -it go-microservice bash

build:
    go build -o main .

deps:
    go mod tidy