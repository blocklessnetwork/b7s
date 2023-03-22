.PHONY: all
all: clean build-node build-keygen

.PHONY: test
test:
	@echo "\n🧪 Testing...\n"
	go clean -testcache
	go test ./src/...
	@echo "\n✅ Done.\n"

.PHONY: build-node
build-node:
	@echo "\n🛠 Building node...\n"
	cd cmd/node && go build -o ../../dist/b7s
	@echo "\n✅ Done.\n"

.PHONY: build-keygen
build-keygen:
	@echo "\n🛠 Building node...\n"
	cd cmd/keygen && go build -o ../../dist/b7s-keygen
	@echo "\n✅ Done.\n"

.PHONY: clean
clean:
	@echo "\n🧹 Cleaning...\n"
	rm -rf dist
	@echo "\n✅ Done.\n"
