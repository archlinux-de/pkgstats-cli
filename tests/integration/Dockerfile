FROM archlinux

RUN pacman -Syu --noconfirm --cachedir /tmp/pacman-cache go just git gcc

COPY . /app/
WORKDIR /app
ENV INTEGRATION_TEST=1
RUN just prepare
RUN just build
RUN just test
RUN just install
RUN systemctl is-enabled -q pkgstats.timer

ENTRYPOINT ["/usr/bin/pkgstats"]
