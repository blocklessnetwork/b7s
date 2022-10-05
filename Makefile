.phony: all
all: clean build

.phony: test
test:
	@echo "Testing..."
	go clean -testcache
	go test ./src/...
	@echo "Done."

.phony: build
build:
	@echo "Building node..."
	cd src && go build -o ../dist/b7s
	@echo "Done."

.phony: clean
clean:
	@echo "Cleaning..."
	rm -rf dist
	@echo "Done."