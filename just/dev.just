mod aarch64
mod arm
mod i686
mod loongarch64
mod riscv64
mod x86_64

# update go modules
update:
    sed -E '/^go\s+[0-9\.]+$/d' -i go.mod
    go get -u -t all
    go mod tidy
    sed -E 's/^go ([0-9]+\.[0-9]+)\.[0-9]+/go \1/' -i go.mod

# run go vet
check-vet:
    go vet ./...

# run static code checks
check-static:
    staticcheck ./...

# run golangci-lint
golangci-lint *args='run':
    golangci-lint {{ args }}

# run gofumpt
gofumpt *args='-w .':
    gofumpt {{ args }}

# check go format
check-fmt:
    test -z $(gofmt -l .)

# run all static checks
check: check-fmt check-vet check-static

# run unit tests on different CPU architectures
test-cross-platform:
    @for a in aarch64 arm i686 loongarch64 riscv64 x86_64; do echo -n "$a "; just -v $a test; done

# build for different CPU architectures
test-build:
    @for a in aarch64 arm i686 loongarch64 riscv64 x86_64; do echo -n "$a "; just -v $a build; done

# test cpu architecture detection on different CPUs
test-cpu-detection:
    @for a in aarch64 arm i686 loongarch64 riscv64 x86_64; do echo -n "$a "; just -v $a test-cpu-detection; done

# test os architecture detection on different CPUs
test-os-detection:
    @for a in aarch64 arm i686 loongarch64 riscv64 x86_64; do echo -n "$a "; just -v $a test-os-detection; done

# run integration tests with a mocked API server
test-integration:
    docker buildx build --pull . -f tests/integration/Dockerfile -t pkgstats-test-integration

# run all available tests
test-all: check test test-cross-platform test-build test-cpu-detection test-os-detection test-integration

# remove any untracked and generated files
clean:
    git clean -fdqx -e .idea
