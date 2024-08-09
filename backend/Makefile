BUILD_DIR="./.bin/api"

build:
	@go build -o "${BUILD_DIR}"

run: build
	@"${BUILD_DIR}"

test:
	@go test -v ./... -count=1

test-only:
	@go test -v ./... --run "$1" -count=1

seed:
	@go run ./scripts/seed.go