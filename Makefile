.PHONY: build up down clean test

# Build the CLI and Server
build:
	@echo "🏗️ Building Lumina-Plane..."
	cd cmd/lumina && go build -o ../../bin/lumina main.go
	cd cmd/server && go build -o ../../bin/server main.go
	@echo "✅ Build complete. Binaries are in /bin"

# Launch the full platform
up:
	@echo "🚀 Launching Lumina-Plane Stack..."
	podman-compose up -d
	@echo "✅ Stack is up! API: http://localhost:8000, Grafana: http://localhost:3000"

# Stop the platform
down:
	@echo "🛑 Stopping Lumina-Plane..."
	podman-compose down
	@echo "✅ Stack stopped."

# Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	rm -rf bin/
	podman-compose down -v
	@echo "✅ Cleaned."

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...
