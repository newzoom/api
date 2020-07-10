.PHONY: test dev remove-infras init
SRC_PATH=$(GOPATH)/src/github.com/newzoom/api
POSTGRES_CONTAINER?=newzoom_db

dev:
	@GO111MODULE=on RUN_MODE=local go run cmd/*.go

build:
	cp -r $(GOPATH)/src/github.com/newzoom/websdk ./bin/assets
	@GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o bin/server cmd/*.go
	docker build --rm -t phuwn29/newzoom .
	rm -rf ./bin/assets
	rm ./bin/server

test:
	go test -p 1 ./pkg/...

remove-infras:
	docker-compose stop; docker-compose  rm -f

seed-db-local:
	@docker cp data/seed/. $(POSTGRES_CONTAINER):/
	@docker exec -t $(POSTGRES_CONTAINER) sh -c "chmod +x seed.sh;./seed.sh"

init: remove-infras
	@docker-compose  up -d 
	@echo "Waiting for database connection..."
	@while ! docker exec $(POSTGRES_CONTAINER) pg_isready -h localhost -p 5432 > /dev/null; do \
		sleep 1; \
	done
	sql-migrate up -config=dbconfig.yml -env="local"
	make seed-db-local

migrate-up-local:
	sql-migrate up -config=dbconfig.yml -env="local"

migrate-down-local:
	sql-migrate down -config=dbconfig.yml -env="local"

migrate-new-local:
	sql-migrate new -config=dbconfig.yml -env="local" $(name)