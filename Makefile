SHELL := /bin/sh

.PHONY: fmt fmt-check lint test ci

fmt:
	@gofmt -s -w .

fmt-check:
	@files="$$(gofmt -s -l .)"; \
	if [ -n "$$files" ]; then \
		echo "gofmt -s found unformatted files:"; echo "$$files"; \
		echo "Run 'make fmt' to format"; \
		exit 1; \
	fi

lint:
	@go vet ./...

test:
	@go test -race -count=1 ./...

ci: fmt-check lint test

