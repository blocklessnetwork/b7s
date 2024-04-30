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
	@UNAME_S=$$(uname -s); \
	UNAME_M=$$(uname -m); \
	if [ "$$UNAME_S" = "Darwin" -a "$$UNAME_M" = "arm64" ]; then \
	    echo "Detected MacOS (arm64). Downloading appropriate version..."; \
	    wget -O /tmp/blockless-runtime.tar.gz https://github.com/blocklessnetwork/runtime/releases/download/v0.3.2/blockless-runtime.macos-latest.aarch64.tar.gz; \
	elif [ "$$UNAME_S" = "Darwin" -a "$$UNAME_M" = "x86_64" ]; then \
	    echo "Detected MacOS (x86_64). Downloading appropriate version..."; \
	    wget -O /tmp/blockless-runtime.tar.gz https://github.com/blocklessnetwork/runtime/releases/download/v0.3.2/blockless-runtime.macos-latest.x86_64.tar.gz; \
	elif [ "$$UNAME_S" = "Linux" -a "$$UNAME_M" = "arm64" ]; then \
	    echo "Detected Linux (arm64). Downloading appropriate version..."; \
	    wget -O /tmp/blockless-runtime.tar.gz https://github.com/blocklessnetwork/runtime/releases/download/v0.3.2/blockless-runtime.linux-latest.arm64.tar.gz; \
	elif [ "$$UNAME_S" = "Linux" -a "$$UNAME_M" = "x86_64" ]; then \
	    echo "Detected Linux (x86_64). Downloading appropriate version..."; \
	    wget -O /tmp/blockless-runtime.tar.gz https://github.com/blocklessnetwork/runtime/releases/download/v0.3.2/blockless-runtime.linux-latest.x86_64.tar.gz; \
	else \
	    echo "No compatible runtime found. Please check your OS and architecture."; \
	fi
	tar -xzf /tmp/blockless-runtime.tar.gz -C /tmp/runtime
	@echo "\n✅ Done.\n"


.PHONY: run-head
run-head:
	@echo "\n🚀 Launching Head Node...\n"
	./dist/b7s --peer-db /tmp/b7s/head-peer-db \
	--function-db /tmp/b7s/head-fdb \
	--log-level debug \
	--port 9527 \
	--role head \
	--workspace /tmp/debug/head \
	--private-key ./configs/testkeys/ident1/priv.bin \
	--rest-api :8081
	@echo "\n✅ Head Node is running!\n"

.PHONY: run-worker
run-worker:
	@echo "\n🚀 Launching Worker Node...\n"
	./dist/b7s --peer-db /tmp/b7s/worker-peer-db \
	--function-db /tmp/b7s/worker-fdb \
	--log-level debug \
	--port 0 \
	--role worker \
	--runtime-path /tmp/runtime \
	--runtime-cli bls-runtime \
	--workspace /tmp/debug/worker \
	--private-key ./configs/testkeys/ident2/priv.bin \
	--boot-nodes /ip4/0.0.0.0/tcp/9527/p2p/12D3KooWH9GerdSEroL2nqjpd2GuE5dwmqNi7uHX7FoywBdKcP4q
	@echo "\n✅ Worker Node is running!\n"
