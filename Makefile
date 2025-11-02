.PHONY: build run clean test deps install

# Build the application
build:
	go build -o bin/smsir cmd/smsir/main.go

# Run the application
run: build
	./bin/smsir

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Install the application
install: build
	go install ./cmd/smsir

# Run tests
test:
	go test ./...

# Run with startup animation
startup: build
	./bin/smsir startup

# Show help
help:
	@echo "Available commands:"
	@echo "  build     - Build the application"
	@echo "  run       - Build and run the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Download and tidy dependencies"
	@echo "  install   - Install the application"
	@echo "  test      - Run tests"
	@echo "  startup   - Run with startup animation"
	@echo "  help      - Show this help"
