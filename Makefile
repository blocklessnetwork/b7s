.PHONY: all
all: clean build

.PHONY: test
test:
	@echo "\n🧪 Testing...\n"
	go clean -testcache
	go test ./src/...
	@echo "\n✅ Done.\n"

.PHONY: build
build:
	@echo "\n🛠 Building node...\n"
	cd src && go build -o ../dist/b7s
	@echo "\n✅ Done.\n"

.PHONY: clean
clean:
	@echo "\n🧹 Cleaning...\n"
	rm -rf dist
	@echo "\n✅ Done.\n"
