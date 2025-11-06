mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
envfile := ./.env
clear_db_after_schema_change := database/last-cleared.dummy
db_schema := database/init/*

.PHONY: help start start-no-webhooks debug sql logs stop clear-db go-build go-run go-test go-air go-test-frontend

# help target adapted from https://gist.github.com/prwhite/8168133#gistcomment-2278355
TARGET_MAX_CHAR_NUM=20

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z_0-9-]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  %-$(TARGET_MAX_CHAR_NUM)s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## Start the services
start: $(envfile) $(clear_db_after_schema_change)
	@echo "Pulling images from Docker Hub (this may take a few minutes)"
	docker compose pull
	@echo "Starting Docker services"
	docker compose up --build --detach
	./wait-for-client.sh

## Start the services without webhooks
start-no-webhooks: $(envfile) $(clear_db_after_schema_change)
	@echo "Pulling images from Docker Hub (this may take a few minutes)"
	docker compose pull
	@echo "Starting Docker services"
	docker compose up --detach client
	./wait-for-client.sh

## Start the services in debug mode
debug: $(envfile) $(clear_db_after_schema_change)
	@echo "Starting services (this may take a few minutes if there are any changes)"
	docker compose -f docker-compose.yml -f docker-compose.debug.yml up --build --detach
	./wait-for-client.sh

## Start an interactive psql session (services must running)
sql:
	docker compose exec db psql -U postgres

## Show the service logs (services must be running)
logs:
	docker compose logs --follow

## Stop the services
stop:
	docker compose down
	docker volume rm $(current_dir)_{client,server}_node_modules 2>/dev/null || true

## Clear the sandbox and production databases
clear-db: stop
	docker volume rm $(current_dir)_pg_{sandbox,production}_data 2>/dev/null || true

## Build the Go server
go-build:
	cd go-server && go build -o server ./cmd/server

## Run the Go server (requires .env and database running)
go-run:
	cd go-server && go run ./cmd/server

## Test the Go server (requires .env and database running)
go-test:
	cd go-server && go test ./...

## Watch and auto-rebuild Go server on file changes with Air (hot-reload)
go-air:
	cd go-server && ~/go/bin/air -c .air.toml

## Test Plaid Link flow: launches Go server + React frontend (hit Ctrl+C to stop both)
go-test-frontend:
	@echo "Starting Go server (port 8000) and React frontend (port 5173)..."
	@echo "React frontend: http://localhost:5173"
	@echo ""
	@echo "Press Ctrl+C to stop both servers"
	@echo ""
	@(cd go-server && ~/go/bin/air -c .air.toml &) && (cd client-go && npm run dev); kill %1 2>/dev/null || true

$(envfile):
	@echo "Error: .env file does not exist! See the README for instructions."
	@exit 1

# Remove local DBs if the DB schema has changed
$(clear_db_after_schema_change): $(db_schema)
	@$(MAKE) clear-db
	@touch $(clear_db_after_schema_change)
