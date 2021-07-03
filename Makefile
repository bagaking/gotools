MODULES := $(shell find . \( -type d \( -name .git -o -name vendor -o -name node_modules \) -prune \) -o \( -name go.mod -exec dirname {} \; \) | sort)

.PHONY: check test
check:
	@set -e; \
	for module in $(MODULES); do \
		echo "==> go test ./... ($$module)"; \
		(cd "$$module" && go test ./...); \
	done

test: check
