APP=sellerpayout
MOCK_FOLDER=${PWD}/pkg/mock
D_PATH=Dockerfile
IGNORED_FOLDER=.ignore
COVERAGE_FILE=$(IGNORED_FOLDER)/coverage.out

DATABASE_DSN=postgres://u:p@localhost:5432/postgres?sslmode=disable

.PHONY: up dev down tools cover-html cover clean test lint mock

##
## Local stack development.
## Do not run following commands in CI
##

up:
	@D_PATH=$(D_PATH) docker-compose up --remove-orphans --build -d
	@docker-compose logs -f ${APP}

##	up local stack in development mode
##	a filewatcher is present for auto-reload the app
dev: D_PATH=Dockerfile.dev
dev: up

# down the local stack
down:
	@docker-compose down


##
## Database
##	

migrate-up: ## Apply migrations not yet done
	@migrate -database ${DATABASE_DSN} -path migrations up

migrate-down: ## Apply a migration down to return to previous database state
	@migrate -database ${DATABASE_DSN} -path migrations down 1

reset-db: ## Drop all database and apply all migrations
	@migrate -database ${DATABASE_DSN} -path migrations drop
	@migrate -database ${DATABASE_DSN} -path migrations up

.PHONY: migrate-up migrate-down reset-db

##
## Quality Code
##

mock:
	@MOCK_FOLDER=${MOCK_FOLDER} go generate ./...

lint:
	@staticcheck ./...

test: mock
	@mkdir -p ${IGNORED_FOLDER}
	@go test -gcflags=-l -count=1 -race -coverprofile=${COVERAGE_FILE} -covermode=atomic ./...

cover: ## Cover
	@if [ ! -e ${COVERAGE_FILE} ]; then \
		echo "Error: ${COVERAGE_FILE} doesn't exists. Please run \`make test\` then retry."; \
		exit 1; \
	fi
	@go tool cover -func=${COVERAGE_FILE}

cover-html: ## Cover html
	@if [ ! -e ${COVERAGE_FILE} ]; then \
		echo "Error: ${COVERAGE_FILE} doesn't exists. Please run \`make test\` then retry."; \
		exit 1; \
	fi
	@go tool cover -html=${COVERAGE_FILE}

clean:
	@rm -rf ${IGNORED_FOLDER}
	@rm -rf ${COVERAGE_FILE}

##
## Docs
##

swag:
	@swag init --parseDependency --parseDepth=2 -g ./internal/handler/server.go -o ./docs

.PHONY: swag

##
## Tooling
##
tools:
	@go install github.com/golang/mock/mockgen@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@GO111MODULE=off go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate