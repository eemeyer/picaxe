GO_SOURCE_FILES=$(shell find . -type f -name "*.go")
BUILD_DIR=build

default: $(BUILD_DIR)/picaxe

clean:
	rm -f $(BUILD_DIR)/picaxe

$(BUILD_DIR)/picaxe: build
	go build -o $(BUILD_DIR)/picaxe .

test: build 
	go test `go list ./... | grep -v /vendor/` | fgrep -v "[no test files]"

build: $(GO_SOURCE_FILES)
	go build `go list ./... | grep -v /vendor/`

.PHONY: default all clean test build picaxe
