.PHONY: db-up db-down api-up api-down run run-seed

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=geolocation
MAIN=cmd/api/main.go

api-up:
	SERVER_MODE=container docker-compose up --remove-orphans --build api
	docker-compose down api
run:
	$(GOCMD) run ./cmd/api
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN)
test:
	$(GOTEST) -v ./...