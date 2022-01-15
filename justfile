export CGO_CPPFLAGS := env_var_or_default('CPPFLAGS', '')
export CGO_CFLAGS := env_var_or_default('CFLAGS', '')
export CGO_CXXFLAGS := env_var_or_default('CXXFLAGS', '')
export CGO_LDFLAGS := env_var_or_default('LDFLAGS', '')

default:
	@just --list

build:
	go build -a -o pkgstats \
		-trimpath -buildmode=pie -mod=readonly -modcacherw \
		-ldflags '-linkmode=external -s -w -X pkgstats-cli/internal/build.Version={{`git describe --tags`}}'

test:
	go vet
	go test -v ./...

test-cross-platform:
	CGO_ENABLED=0 GOARCH=arm go test -v -exec qemu-arm ./...
	CGO_ENABLED=0 GOARCH=arm64 go test -v -exec qemu-aarch64 ./...
	CGO_ENABLED=0 GOARCH=riscv64 go test -v -exec qemu-riscv64 ./...
	CGO_ENABLED=0 GOARCH=386 go test -v -exec linux32 ./...

test-build:
	@for arch in amd64 386 arm64 arm riscv64; do \
		echo "Building for ${arch}"; \
		CGO_ENABLED=0 GOARCH=${arch} go build -o tests/build/pkgstats-build-${arch}; \
	done

test-cpu-detection: test-build
	@# ARM 32-Bit
	qemu-arm -cpu arm946 ./tests/build/pkgstats-build-arm submit --dump-json | jq -r '.system.architecture' | grep -q '^armv5$'
	qemu-arm -cpu arm1176 ./tests/build/pkgstats-build-arm submit --dump-json | jq -r '.system.architecture' | grep -q '^armv6$'
	qemu-arm -cpu cortex-a15 ./tests/build/pkgstats-build-arm submit --dump-json | jq -r '.system.architecture' | grep -q '^armv7$'
	qemu-arm -cpu max ./tests/build/pkgstats-build-arm submit --dump-json | jq -r '.system.architecture' | grep -q '^aarch64$'
	@# ARM 64-Bit
	qemu-aarch64 ./tests/build/pkgstats-build-arm64 submit --dump-json | jq -r '.system.architecture' | grep -q '^aarch64$'
	@# RISC-V 64-Bit rv64gc
	qemu-riscv64 -cpu sifive-u54 ./tests/build/pkgstats-build-riscv64 submit --dump-json | jq -r '.system.architecture' | grep -q '^riscv64$'
	@# x86_64
	qemu-x86_64 -cpu Conroe ./tests/build/pkgstats-build-amd64 submit --dump-json | jq -r '.system.architecture' | grep -q '^x86_64$'
	qemu-x86_64 -cpu Nehalem ./tests/build/pkgstats-build-amd64 submit --dump-json | jq -r '.system.architecture' | grep -q '^x86_64_v2$'
	@# 32-Bit on x86_64
	linux32 ./tests/build/pkgstats-build-386 submit --dump-json | jq -r '.system.architecture' | grep -q '^x86_64'

test-integration:
	docker build --pull . -f tests/integration/Dockerfile -t pkgstats-test-integration

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

test-all: test test-build test-cpu-detection test-integration

clean:
	git clean -fdqx -e .idea

# vim: set ft=make :
