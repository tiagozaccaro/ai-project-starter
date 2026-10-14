.PHONY: all build test vet fmt lint clean

BINARY := project-forge
CMD := ./cmd/project-forge

all: build test vet

build:
	go build -o $(BINARY) $(CMD)

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

lint: vet
	@echo "Checking formatting..."
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not formatted:"; \
		echo "$$unformatted"; \
		echo "Run 'make fmt' to fix."; \
		exit 1; \
	fi

clean:
	rm -f $(BINARY)

run:
	go run $(CMD)

coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total

smoke:
	./tests/smoke.sh

smoke-fullstack:
	./tests/smoke.sh fullstack

smoke-backend:
	./tests/smoke.sh backend

smoke-frontend:
	./tests/smoke.sh frontend
