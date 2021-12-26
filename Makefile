.PHONY: all build test test-build test-cpu-detection test-integration install clean

all: build

VERSION != git describe --tags

export CGO_CPPFLAGS=${CPPFLAGS}
export CGO_CFLAGS=${CFLAGS}
export CGO_CXXFLAGS=${CXXFLAGS}
export CGO_LDFLAGS=${LDFLAGS}

build:
	go build -a -o pkgstats -trimpath -buildmode=pie -mod=readonly -modcacherw -ldflags '-s -w -X pkgstats-cli/internal/build.Version=${VERSION}'

test:
	go vet
	go test -v ./...

test-build:
	@for arch in amd64 386 arm64 arm; do \
		echo "Building for $${arch}"; \
		CGO_ENABLED=0 GOARCH=$${arch} go build -o pkgstats-build-$${arch}; \
	done

test-cpu-detection: test-build
	@# ARM 32-Bit
	qemu-arm -cpu arm946 ./pkgstats-build-arm submit --dump-json | jq -r '.system.architecture' | grep -q '^armv5$$'
	qemu-arm -cpu arm1176 ./pkgstats-build-arm submit --dump-json | jq -r '.system.architecture' | grep -q '^armv6$$'
	qemu-arm -cpu cortex-a15 ./pkgstats-build-arm submit --dump-json | jq -r '.system.architecture' | grep -q '^armv7$$'
	qemu-arm -cpu max ./pkgstats-build-arm submit --dump-json | jq -r '.system.architecture' | grep -q '^aarch64$$'
	@# ARM 64-Bit
	qemu-aarch64 ./pkgstats-build-arm64 submit --dump-json | jq -r '.system.architecture' | grep -q '^aarch64$$'
	@# x86_64
	qemu-x86_64 -cpu Conroe pkgstats-build-amd64 submit --dump-json | jq -r '.system.architecture' | grep -q '^x86_64$$'
	qemu-x86_64 -cpu Nehalem pkgstats-build-amd64 submit --dump-json | jq -r '.system.architecture' | grep -q '^x86_64_v2$$'

test-integration:
	docker build --pull . -t pkgstats

install:
	# cli
	install -D pkgstats -m755 "$(DESTDIR)/usr/bin/pkgstats"

	# systemd timer
	for service in pkgstats.service pkgstats.timer; do \
		install -Dt "$(DESTDIR)/usr/lib/systemd/system" -m644 init/$${service} ; \
	done
	install -d "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants"
	cd "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants" && ln -s ../pkgstats.timer

	# bash completions
	install -d "$(DESTDIR)/usr/share/bash-completion/completions"
	./pkgstats completion bash > "$(DESTDIR)/usr/share/bash-completion/completions/pkgstats"

	# zsh completions
	install -d "$(DESTDIR)/usr/share/zsh/site-functions/"
	./pkgstats completion zsh > "$(DESTDIR)/usr/share/zsh/site-functions/_pkgstats"

	# fish completions
	install -d "$(DESTDIR)/usr/share/fish/vendor_completions.d"
	./pkgstats completion fish > "$(DESTDIR)/usr/share/fish/vendor_completions.d/pkgstats.fish"

clean:
	git clean -fdqx -e .idea
