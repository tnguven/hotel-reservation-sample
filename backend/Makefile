BUILD_DIR="./.bin/api"
GO_FLAGS=GOEXPERIMENT=rangefunc

build:
	@${GO_FLAGS} go build -o "${BUILD_DIR}"

run: build
	@${GO_FLAGS} "${BUILD_DIR}"

# test:
# 	@go test -v ./... -count=1

test-rerun:
	@${GO_FLAGS} go run gotest.tools/gotestsum@latest --debug --format testname --rerun-fails --packages="./..." -- -count=2

test:
	@${GO_FLAGS} go run gotest.tools/gotestsum@latest --debug --format testname --debug --packages="./..." -- -count=1

test-watch:
	@${GO_FLAGS} go run gotest.tools/gotestsum@latest --debug --format testname --watch --packages="./..." -- -count=1

test-only:
	@${GO_FLAGS} go run gotest.tools/gotestsum@latest --debug --format standard-verbose --watch --packages="./..." -- --run $(TEST_NAME) -count=1

seed:
	@${GO_FLAGS} go run ./scripts/seed.go