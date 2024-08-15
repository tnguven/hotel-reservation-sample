BUILD_DIR="./.bin/api"

build:
	@go build -o "${BUILD_DIR}"

run: build
	@"${BUILD_DIR}"

# test:
# 	@go test -v ./... -count=1

test-rerun:
	@GOEXPERIMENT=rangefunc go run gotest.tools/gotestsum@latest --debug --format testname --rerun-fails --packages="./..." -- -count=2

test:
	@GOEXPERIMENT=rangefunc go run gotest.tools/gotestsum@latest --debug --format testname --debug --packages="./..." -- -count=1

test-watch:
	@GOEXPERIMENT=rangefunc go run gotest.tools/gotestsum@latest --debug --format testname --watch --packages="./..." -- -count=1

test-only:
	@GOEXPERIMENT=rangefunc go run gotest.tools/gotestsum@latest --debug --format standard-verbose --watch --packages="./..." -- --run $(TEST_NAME) -count=1

seed:
	@GOEXPERIMENT=rangefunc go run ./scripts/seed.go