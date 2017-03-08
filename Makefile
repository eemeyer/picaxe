GO_SOURCE_FILES=$(shell find . -type f -name "*.go")
BUILD_DIR=build

default: build

clean:
	rm -f $(BUILD_DIR)/picaxe

build: $(BUILD_DIR)/picaxe

$(BUILD_DIR)/picaxe: $(GO_SOURCE_FILES)
	go build -o $(BUILD_DIR)/picaxe .

test: build 
	go test `go list ./... | grep -v /vendor/` | fgrep -v "[no test files]"

run: build
	$(BUILD_DIR)/picaxe -l localhost

.PHONY: default all clean test run build picaxe
