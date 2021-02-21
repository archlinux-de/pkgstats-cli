FROM alpine

RUN apk add go make git bats php php-json jq pacman
RUN echo -e '[core]\nSigLevel=Never\nServer=https://mirror.rackspace.com/archlinux/$repo/os/$arch' >> /etc/pacman.conf
RUN pacman -Sy --noconfirm --noscriptlet --quiet pacman-mirrorlist archlinux-keyring

COPY . /app/
WORKDIR /app
RUN make build
RUN make test
RUN php -S localhost:8888 tests/server.php & sleep 2 && bats tests/integration.bats
RUN make DESTDIR=/ install

ENTRYPOINT ["/usr/bin/pkgstats"]
