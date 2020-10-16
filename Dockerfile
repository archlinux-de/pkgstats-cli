FROM archlinux
RUN pacman -Syu --noconfirm go make gcc git
COPY . /app/
WORKDIR /app
RUN make build
RUN make DESTDIR=/ install
ENTRYPOINT ["/usr/bin/pkgstats"]
