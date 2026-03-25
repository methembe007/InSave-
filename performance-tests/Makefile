.PHONY: help install normal peak stress all clean

# Default target
help:
	@echo "InSavein Performance Testing"
	@echo ""
	@echo "Available targets:"
	@echo "  make install    - Install k6 (Linux/macOS)"
	@echo "  make normal     - Run normal load test"
	@echo "  make peak       - Run peak load test"
	@echo "  make stress     - Run stress test"
	@echo "  make all        - Run all tests sequentially"
	@echo "  make clean      - Clean up result files"
	@echo ""
	@echo "Environment variables:"
	@echo "  BASE_URL        - Base URL for testing (default: http://localhost:8080)"

# Install k6
install:
	@echo "Installing k6..."
	@if command -v brew >/dev/null 2>&1; then \
		brew install k6; \
	elif command -v apt-get >/dev/null 2>&1; then \
		sudo gpg -k && \
		sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69 && \
		echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list && \
		sudo apt-get update && \
		sudo apt-get install k6; \
	else \
		echo "Please install k6 manually from https://k6.io/docs/getting-started/installation/"; \
	fi

# Run normal load test
normal:
	@./run-tests.sh normal

# Run peak load test
peak:
	@./run-tests.sh peak

# Run stress test
stress:
	@./run-tests.sh stress

# Run all tests
all:
	@./run-tests.sh all

# Clean up results
clean:
	@echo "Cleaning up result files..."
	@rm -rf results/
	@echo "Done"
