.PHONY: all build test install

all: build

build:
	sed "s/@VERSION@/$$(git describe --tags)/g" pkgstats.sh > pkgstats

test:
	shellcheck pkgstats.sh
	bats tests

install:
	install -D pkgstats -m755 "$(DESTDIR)/usr/bin/pkgstats"
	install -Dt "$(DESTDIR)/usr/lib/systemd/system" -m644 pkgstats.{timer,service}
	install -d "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants"
	ln -st "$(DESTDIR)/usr/lib/systemd/system/timers.target.wants" ../pkgstats.timer
