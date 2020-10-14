.PHONY: all build test install clean

all: build

VERSION != git describe --tags || echo 2.5.0-0-dev

build:
	go build -o pkgstats -trimpath -buildmode=pie -ldflags '-linkmode external -extldflags "${LDFLAGS}" -X main.Version=${VERSION}'

test:
	go vet
	go test -v

install:
	install -D pkgstats -m755 "$(DESTDIR)/usr/bin/pkgstats"
	install -Dt "$(DESTDIR)/usr/lib/systemd/system" -m644 pkgstats.{timer,service}
	install -d "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants"
	ln -st "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants" ../pkgstats.timer

clean:
	git clean -fdqx -e .idea
