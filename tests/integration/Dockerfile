FROM alpine

RUN apk add go just git bats php php-json jq pacman
RUN echo -e '[core]\nSigLevel=Never\nServer=https://mirror.rackspace.com/archlinux/$repo/os/$arch' >> /etc/pacman.conf
RUN pacman -Sy --noconfirm --noscriptlet --quiet pacman-mirrorlist archlinux-keyring

COPY . /app/
WORKDIR /app
RUN just build
RUN just test
RUN php -S localhost:8888 tests/integration/server.php & sleep 2 && bats tests/integration/integration.bats
RUN just install

ENTRYPOINT ["/usr/bin/pkgstats"]
