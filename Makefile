MODULES := $(shell find . -name go.mod ! -path './.git/*' ! -path './vendor/*' -exec dirname {} \; | sort)

.PHONY: check test
check:
	@set -e; \
	for module in $(MODULES); do \
		echo "==> go test ./... ($$module)"; \
		(cd "$$module" && go test ./...); \
	done

test: check
