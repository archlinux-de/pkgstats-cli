FROM archlinux

RUN pacman -Syu --noconfirm --cachedir /tmp/pacman-cache rust just git bats php jq pkgconf

COPY . /app/
WORKDIR /app
RUN just prepare
RUN just build
RUN just test
RUN php -S localhost:8888 tests/integration/server.php & sleep 2 && bats tests/integration/integration.bats
RUN just install

ENTRYPOINT ["/usr/bin/pkgstats"]
