# COLORS
ccgreen=$(shell echo "\033[32m")
ccred=$(shell echo "\033[0;31m")
ccyellow=$(shell echo "\033[0;33m")
ccend=$(shell echo "\033[0m")

# SILENT MODE (avoid echoes)
.SILENT: all fmt test linter build

# PROCESS
all: fmt test linter build

fmt:
	echo "$(ccyellow)Formatting files...$(ccend)"
	$(GOPATH)/bin/goimports -w -local github.com/alexyslozada/shorturl .
	echo "$(ccgreen)Formatting files done!$(ccend)"

test:
	for d in $$(go list ./...); do \
		if go test -v -failfast $$d; then \
			echo "$(ccyellow)$$d test pass!!!$(ccend)"; \
		else \
			echo "$(ccred)$$d test failed :($(ccend)"; \
			exit 1; \
		fi; \
	done;
	echo "$(ccgreen)All test pass!$(ccend)"

install-linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.46.2

linter:
	echo "$(ccyellow)Executing linter...$(ccend)"
	golangci-lint run
	echo "$(ccyellow)Linter finished!$(ccend)"

build:
	echo "$(ccyellow)Building app...$(ccend)"
	go build -o short ./cmd/echopostgres/main.go
	echo "$(ccgreen)Finish build!$(ccend)"
