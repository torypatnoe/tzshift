# tzshift — build / test / release
BINARY := tzshift
PKG    := ./cmd/tzshift
DIST   := dist
LDFLAGS := -s -w

# Statically linked release targets (no host zoneinfo needed — IANA DB embedded).
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

.PHONY: build test vet fmt install release clean

build: ## Build the binary for the host platform
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(PKG)

test: ## Run the test suite
	go test ./...

vet: ## go vet
	go vet ./...

fmt: ## Check formatting
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed:"; gofmt -l .; exit 1)

install: ## Install to GOBIN/GOPATH
	go install -ldflags "$(LDFLAGS)" $(PKG)

release: ## Cross-compile static binaries into dist/
	@mkdir -p $(DIST)
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; \
		echo "building $$os/$$arch"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch \
			go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-$$os-$$arch $(PKG) || exit 1; \
	done
	@echo "→ $(DIST)/"

clean:
	rm -rf $(DIST) $(BINARY)
