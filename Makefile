.PHONY: all
all: clean build-node build-keyforge build-manager

.PHONY: test
test:
	@echo "\n🧪 Testing...\n"
	go clean -testcache
	go test ./...
	@echo "\n✅ Done.\n"

.PHONY: build-node
build-node:
	@echo "\n🛠 Building node...\n"
	cd cmd/node && go build -o ../../dist/b7s
	@echo "\n✅ Done.\n"

.PHONY: build-keyforge
build-keyforge:
	@echo "\n🛠 Building node keygen...\n"
	cd cmd/keyforge && go build -o ../../dist/b7s-keyforge
	@echo "\n✅ Done.\n"

.PHONY: build-manager
build-manager:
	@echo "\n🛠 Building node manager...\n"
	cd cmd/manager && go build -o ../../dist/b7s-manager
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
	wget -O /tmp/blockless-runtime.tar.gz https://github.com/blocklessnetwork/runtime/releases/download/v0.3.4/blockless-runtime.darwin-latest.x86_64.tar.gz
	tar -xzf /tmp/blockless-runtime.tar.gz -C /tmp/runtime
	@echo "\n✅ Done.\n"
