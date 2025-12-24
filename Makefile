APP_NAME := OpenLinkHub-system-tray
SRC_DIR := src
BUILD_DIR := bin

GO := go
GOFLAGS :=

.PHONY: all build run clean tidy

all: build

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./$(SRC_DIR)

run:
	$(GO) run ./$(SRC_DIR)

tidy:
	$(GO) mod tidy

clean:
	rm -rf $(BUILD_DIR)
