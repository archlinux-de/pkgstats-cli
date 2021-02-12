FROM archlinux

# See https://github.com/qutebrowser/qutebrowser/commit/478e4de7bd1f26bebdcdc166d5369b2b5142c3e2
# and https://bugs.archlinux.org/task/69563
# WORKAROUND for glibc 2.33 and old Docker
# See https://github.com/actions/virtual-environments/issues/2658
# Thanks to https://github.com/lxqt/lxqt-panel/pull/1562
RUN patched_glibc=glibc-linux4-2.33-4-x86_64.pkg.tar.zst && \
    curl -LO "https://repo.archlinuxcn.org/x86_64/$patched_glibc" && \
    bsdtar -C / -xvf "$patched_glibc"

RUN pacman -Syu --noconfirm go make gcc git bash-bats php jq
COPY . /app/
WORKDIR /app
RUN . /etc/makepkg.conf && export CPPFLAGS CFLAGS CXXFLAGS LDFLAGS && make build
RUN make test
RUN php -S localhost:8888 tests/server.php & sleep 2 && bats tests/integration.bats
RUN make DESTDIR=/ install
ENTRYPOINT ["/usr/bin/pkgstats"]
