mod aarch64
mod arm
mod i686
mod loongarch64
mod riscv64
mod x86_64

# update go modules
[group('dev')]
update:
    sed -E '/^go\s+[0-9\.]+$/d' -i go.mod
    go get -u -t all
    go mod tidy
    sed -E 's/^go ([0-9]+\.[0-9]+)\.[0-9]+/go \1/' -i go.mod

# run golangci-lint
[group('test')]
lint *args='run':
       golangci-lint {{ args }}

# run golangci-lint fmt
[group('dev')]
fmt *args:
    golangci-lint fmt {{ args }}

# run unit tests on different CPU architectures
[group('test')]
[parallel]
test-cross-platform: aarch64::test arm::test i686::test loongarch64::test riscv64::test x86_64::test

# build for different CPU architectures
[group('test')]
[parallel]
test-build: aarch64::build arm::build i686::build loongarch64::build riscv64::build x86_64::build

# test cpu architecture detection on different CPUs
[group('test')]
[parallel]
test-cpu-detection: aarch64::test-cpu-detection arm::test-cpu-detection i686::test-cpu-detection loongarch64::test-cpu-detection riscv64::test-cpu-detection x86_64::test-cpu-detection

# test os architecture detection on different CPUs
[group('test')]
[parallel]
test-os-detection: aarch64::test-os-detection arm::test-os-detection i686::test-os-detection loongarch64::test-os-detection riscv64::test-os-detection x86_64::test-os-detection

# run integration tests with a mocked API server
[group('test')]
test-integration:
    docker buildx build --pull . -f tests/integration/Dockerfile -t pkgstats-test-integration

# run all available tests
[group('test')]
[parallel]
test-all: lint test test-cross-platform test-build test-cpu-detection test-os-detection test-integration

# remove any untracked and generated files
[group('dev')]
clean:
    git clean -fdqx -e .idea

# calculate test coverage and open in a browsertest
[group('dev')]
coverage:
    #!/usr/bin/env bash
    set -euo pipefail
    COVER_FILE=$(mktemp)
    go test -coverprofile "$COVER_FILE" ./...
    go tool cover -html="$COVER_FILE"
    rm -f "$COVER_FILE"
