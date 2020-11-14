.PHONY: all build test test-integration install clean

all: build

VERSION != git describe --tags

export CGO_CPPFLAGS=${CPPFLAGS}
export CGO_CFLAGS=${CFLAGS}
export CGO_CXXFLAGS=${CXXFLAGS}

build:
	go build -a -o pkgstats -trimpath -buildmode=pie -mod=readonly -modcacherw -ldflags '-s -w -X pkgstats-cli/internal/build.Version=${VERSION}'

test:
	go vet
	go test -v ./...

test-integration:
	docker build --pull . -t pkgstats

install:
	# cli
	install -D pkgstats -m755 "$(DESTDIR)/usr/bin/pkgstats"

	# systemd timer
	install -Dt "$(DESTDIR)/usr/lib/systemd/system" -m644 init/pkgstats.{timer,service}
	install -d "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants"
	ln -st "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants" ../pkgstats.timer

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
