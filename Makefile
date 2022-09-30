build-node:
	@echo "Building node..."
	cd src && go build -o ../dist/b7s
	@echo "Done."