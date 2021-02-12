FROM archlinux
RUN mkdir -p /var/lib/pacman
RUN pacman -Syu --noconfirm go make gcc git bash-bats php jq
COPY . /app/
WORKDIR /app
RUN . /etc/makepkg.conf && export CPPFLAGS CFLAGS CXXFLAGS LDFLAGS && make build
RUN make test
RUN php -S localhost:8888 tests/server.php & sleep 2 && bats tests/integration.bats
RUN make DESTDIR=/ install
ENTRYPOINT ["/usr/bin/pkgstats"]
