MODULES := . annotation csvp fuctx procast reflectool strs

.PHONY: check test
check:
	@set -e; \
	for module in $(MODULES); do \
		echo "==> go test ./... ($$module)"; \
		(cd "$$module" && go test ./...); \
	done

test: check
