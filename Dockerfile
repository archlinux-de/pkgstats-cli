FROM archlinux
RUN pacman -Syu --noconfirm go make gcc git bash-bats
COPY . /app/
WORKDIR /app
RUN make build
RUN make test
RUN bats integration.bats
RUN make DESTDIR=/ install
ENTRYPOINT ["/usr/bin/pkgstats"]
