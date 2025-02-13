BUILD_DIR="./.bin/svc-api"
GO_FLAGS=GOEXPERIMENT=rangefunc

build:
	@${GO_FLAGS} go build -o .bin/svc-api cmd/svc-api/

run: build
	@${GO_FLAGS} "${BUILD_DIR}"

test-rerun:
	@${GO_FLAGS} go run gotest.tools/gotestsum@latest --debug --format testname --rerun-fails --packages="./..." -- -count=2

test:
	@${GO_FLAGS} go run gotest.tools/gotestsum@latest --debug --format testname --debug --packages="./..." -- -count=1

test-watch:
	@${GO_FLAGS} go run gotest.tools/gotestsum@latest --debug --format testname --watch --packages="./..." -- -count=1

test-only:
	@${GO_FLAGS} go run gotest.tools/gotestsum@latest --debug --format standard-verbose --watch --packages="./..." -- --run $(TEST_NAME) -count=1

seed:
	@${GO_FLAGS} go run ./cmd/task-seeder
