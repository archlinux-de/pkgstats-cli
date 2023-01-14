export CGO_CPPFLAGS := env_var_or_default('CPPFLAGS', '')
export CGO_CFLAGS := env_var_or_default('CFLAGS', '')
export CGO_CXXFLAGS := env_var_or_default('CXXFLAGS', '')
export CGO_LDFLAGS := env_var_or_default('LDFLAGS', '')

# list all recipes
default:
	@just --list

# build pkgstats for production
build:
	go build -a -o pkgstats \
		-trimpath -buildmode=pie -mod=readonly -modcacherw \
		-ldflags '-linkmode=external -s -w -X pkgstats-cli/internal/build.Version={{`git describe --tags`}}'

# update go modules
update:
	sed -E '/^go\s+[0-9\.]+$/d' -i go.mod
	go get -u -t all
	go mod tidy

# run go vet
check-vet:
	go vet ./...

# run static code checks
check-static:
	staticcheck ./...

# check go format
check-fmt:
	test -z $(gofmt -l .)

# run all static checks
check: check-fmt check-vet check-static

# run unit tests
test:
	go test -v ./...

# run unit tests on different CPU architectures
test-cross-platform:
	CGO_ENABLED=0 GOARCH=arm go test -v -exec qemu-arm ./...
	CGO_ENABLED=0 GOARCH=arm64 go test -v -exec qemu-aarch64 ./...
	CGO_ENABLED=0 GOARCH=riscv64 go test -v -exec qemu-riscv64 ./...
	CGO_ENABLED=0 GOARCH=386 go test -v -exec linux32 ./...

# build for different CPU architectures
test-build:
	@for arch in amd64 386 arm64 arm riscv64; do \
		echo "Building for ${arch}"; \
		CGO_ENABLED=0 GOARCH=${arch} go build -o tests/build/pkgstats-${arch}; \
	done

# test cpu architecture detection on different CPUs
test-cpu-detection:
	@# ARM 32-Bit
	CGO_ENABLED=0 GOARCH=arm go run -exec 'qemu-arm -cpu arm946' main.go architecture system | grep -q '^armv5$'
	CGO_ENABLED=0 GOARCH=arm go run -exec 'qemu-arm -cpu arm1176' main.go architecture system | grep -q '^armv6$'
	CGO_ENABLED=0 GOARCH=arm go run -exec 'qemu-arm -cpu cortex-a15' main.go architecture system | grep -q '^armv7$'
	CGO_ENABLED=0 GOARCH=arm go run -exec 'qemu-arm -cpu max' main.go architecture system | grep -q '^aarch64$'
	@# ARM 64-Bit
	CGO_ENABLED=0 GOARCH=arm64 go run -exec 'qemu-aarch64' main.go architecture system | grep -q '^aarch64$'
	@# RISC-V 64-Bit rv64gc
	CGO_ENABLED=0 GOARCH=riscv64 go run -exec 'qemu-riscv64 -cpu sifive-u54' main.go architecture system | grep -q '^riscv64$'
	@# x86_64
	CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64 -cpu Conroe' main.go architecture system | grep -q '^x86_64$'
	CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64 -cpu Nehalem' main.go architecture system | grep -q '^x86_64_v2$'
	@# Test crashes on older Qemu versions
	if qemu-x86_64 -version | grep -Eq 'version [7-9]\.[2-9][0-9]*\.[0-9]+$'; then CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64 -cpu Haswell' main.go architecture system 2>&1 | grep -q '^x86_64_v3$'; fi
	@# 32-Bit on x86_64
	CGO_ENABLED=0 GOARCH=386 go run -exec 'linux32' main.go architecture system | grep -q '^x86_64'

# test os architecture detection on different CPUs
test-os-detection:
	@# ARM 32-Bit
	CGO_ENABLED=0 GOARCH=arm go run -exec 'qemu-arm' main.go architecture os | grep -q '^armv7l$'
	@# ARM 64-Bit
	CGO_ENABLED=0 GOARCH=arm64 go run -exec 'qemu-aarch64' main.go architecture os | grep -q '^aarch64$'
	@# RISC-V 64-Bit rv64gc
	CGO_ENABLED=0 GOARCH=riscv64 go run -exec 'qemu-riscv64' main.go architecture os | grep -q '^riscv64$'
	@# x86_64
	CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64' main.go architecture os | grep -q '^x86_64$'
	@# 32-Bit on x86_64
	CGO_ENABLED=0 GOARCH=386 go run -exec 'linux32' main.go architecture os | grep -q '^i686$'

# run integration tests with a mocked API server
test-integration:
	docker build --pull . -f tests/integration/Dockerfile -t pkgstats-test-integration

# install pkgstats and its configuration
install *DESTDIR='':
	@# cli
	install -D pkgstats -m755 "{{DESTDIR}}/usr/bin/pkgstats"

	@# systemd timer
	for service in pkgstats.service pkgstats.timer; do \
		install -Dt "{{DESTDIR}}/usr/lib/systemd/system" -m644 init/${service} ; \
	done
	install -d "{{DESTDIR}}/usr/lib/systemd/system/timers.target.wants"
	cd "{{DESTDIR}}/usr/lib/systemd/system/timers.target.wants" && ln -s ../pkgstats.timer

	@# bash completions
	install -d "{{DESTDIR}}/usr/share/bash-completion/completions"
	./pkgstats completion bash > "{{DESTDIR}}/usr/share/bash-completion/completions/pkgstats"

	@# zsh completions
	install -d "{{DESTDIR}}/usr/share/zsh/site-functions/"
	./pkgstats completion zsh > "{{DESTDIR}}/usr/share/zsh/site-functions/_pkgstats"

	@# fish completions
	install -d "{{DESTDIR}}/usr/share/fish/vendor_completions.d"
	./pkgstats completion fish > "{{DESTDIR}}/usr/share/fish/vendor_completions.d/pkgstats.fish"

# run all available tests
test-all: check test test-build test-cpu-detection test-os-detection test-integration

# remove any untracked and generated files
clean:
	git clean -fdqx -e .idea

# vim: set ft=make :
