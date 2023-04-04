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

.PHONY: setup
setup:
	@echo "\n📥 Downloading and extracting runtime...\n"
	mkdir -p /tmp/runtime
	wget -O /tmp/blockless-runtime.tar.gz https://github.com/blocklessnetwork/runtime/releases/download/v0.0.12/blockless-runtime.linux-latest.x86_64.tar.gz
	tar -xzf /tmp/blockless-runtime.tar.gz -C /tmp/runtime
	@echo "\n✅ Done.\n"
