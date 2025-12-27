GOVERSION ?= $(shell go env GOVERSION)

.PHONY: audit
audit:
	go mod verify
	go vet ./...
	GOTOOLCHAIN=$(GOVERSION) go run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: format
format:
	GOTOOLCHAIN=$(GOVERSION) go run mvdan.cc/gofumpt@latest -w -l .

.PHONY: lint
lint:
	GOTOOLCHAIN=$(GOVERSION) go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0 run ./...

.PHONY: modernize
modernize:
	GOTOOLCHAIN=$(GOVERSION) go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test=false ./...

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: run-client
run-client:
	go run cmd/client/main.go

.PHONY: run-server
run-server:
	go run cmd/server/main.go