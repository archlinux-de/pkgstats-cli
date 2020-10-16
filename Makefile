.PHONY: all build test test-integration install clean

all: build

VERSION != git describe --tags

export CGO_CPPFLAGS=${CPPFLAGS}
export CGO_CFLAGS=${CFLAGS}
export CGO_CXXFLAGS=${CXXFLAGS}
export CGO_LDFLAGS=${LDFLAGS}

build:
	go build -a -o pkgstats -trimpath -buildmode=pie -mod=readonly -modcacherw -ldflags '-X pkgstats-cli/internal/build.Version=${VERSION} -linkmode external -extldflags "${LDFLAGS}"'

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
