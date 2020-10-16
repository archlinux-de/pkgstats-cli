FROM archlinux
RUN pacman -Syu --noconfirm go make gcc git bash-bats php
COPY . /app/
WORKDIR /app
RUN make build
RUN make test
RUN php -S localhost:8888 server.php&sleep 2&&bats integration.bats
RUN make DESTDIR=/ install
ENTRYPOINT ["/usr/bin/pkgstats"]
