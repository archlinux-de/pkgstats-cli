FROM archlinux
RUN pacman -Syu --noconfirm go make gcc git bash-bats php
COPY . /app/
WORKDIR /app
RUN make build
RUN make test
RUN php -S localhost:8888 tests/server.php & sleep 2 && bats tests/integration.bats
RUN make DESTDIR=/ install
ENTRYPOINT ["/usr/bin/pkgstats"]
