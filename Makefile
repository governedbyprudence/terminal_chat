BINARY  := terminal_chat
BIN_DIR := bin
PKG     := ./...
ROOMS   ?= Lobby,Go,Random
ADDR    ?= :9000
HOST    ?= localhost:9000

GO      ?= go

.DEFAULT_GOAL := help

.PHONY: help build run rooms server client test cover vet fmt lint tidy clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary into ./bin
	$(GO) build -o $(BIN_DIR)/$(BINARY) ./tui

run: build ## Run in standalone mode (single user)
	./$(BIN_DIR)/$(BINARY)

rooms: build ## Run with custom rooms, e.g. make rooms ROOMS="Lobby,Go"
	./$(BIN_DIR)/$(BINARY) --rooms "$(ROOMS)"

server: build ## Start server + TUI, e.g. make server ADDR=:9000
	./$(BIN_DIR)/$(BINARY) --server $(ADDR)

client: build ## Connect to a server, e.g. make client HOST=localhost:9000
	./$(BIN_DIR)/$(BINARY) --connect $(HOST)

test: ## Run tests
	$(GO) test ./tui/...

cover: ## Run tests with coverage
	$(GO) test ./tui/... -cover

vet: ## Run go vet
	$(GO) vet $(PKG)

fmt: ## Format sources
	gofmt -w ./tui

lint: vet ## Run vet + check formatting (read-only)
	@out=$$(gofmt -l ./tui); if [ -n "$$out" ]; then echo "gofmt needed:"; echo "$$out"; exit 1; fi

tidy: ## Tidy module deps
	$(GO) mod tidy

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)
