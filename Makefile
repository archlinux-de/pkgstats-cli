.PHONY: all build test install

all: build

#export CGO_CPPFLAGS="${CPPFLAGS}"
#export CGO_CFLAGS="${CFLAGS}"
#export CGO_CXXFLAGS="${CXXFLAGS}"
#export CGO_LDFLAGS="${LDFLAGS}"
#export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"

build:
	#sed "s/@VERSION@/$$(git describe --tags)/g" pkgstats.sh > pkgstats
	go build -ldflags "-X main.Version=$$(git describe --tags)" -o pkgstats cmd/main.go

test:
	#shellcheck pkgstats.sh
	bats tests

install:
	install -D pkgstats -m755 "$(DESTDIR)/usr/bin/pkgstats"
	install -Dt "$(DESTDIR)/usr/lib/systemd/system" -m644 pkgstats.{timer,service}
	install -d "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants"
	ln -st "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants" ../pkgstats.timer
