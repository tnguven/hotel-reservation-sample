BUILD_DIR="./.bin/api"

build:
	@go build -o "${BUILD_DIR}"

run: build
	@"${BUILD_DIR}"

# test:
# 	@go test -v ./... -count=1

test-rerun:
	@go run gotest.tools/gotestsum@latest --debug --format testname --rerun-fails --packages="./..." -- -count=2

test:
	@go run gotest.tools/gotestsum@latest --debug --format testname --packages="./..." -- -count=1

test-only:
	@go run gotest.tools/gotestsum@latest --debug --format testname --packages="./..." -- --run $(TEST_NAME) -count=1

seed:
	@go run ./scripts/seed.go