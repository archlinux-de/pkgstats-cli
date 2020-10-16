.PHONY: all build test test-integration install clean

all: build

VERSION != git describe --tags

build:
	go build -o pkgstats -trimpath -buildmode=pie -ldflags '-X pkgstats-cli/internal/build.Version=${VERSION}'

test:
	go vet
	go test -v ./...

test-integration:
	docker build . -t pkgstats

install:
	install -D pkgstats -m755 "$(DESTDIR)/usr/bin/pkgstats"
	install -Dt "$(DESTDIR)/usr/lib/systemd/system" -m644 init/pkgstats.{timer,service}
	install -d "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants"
	ln -st "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants" ../pkgstats.timer

clean:
	git clean -fdqx -e .idea
