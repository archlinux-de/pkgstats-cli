FROM archlinux
RUN pacman -Sy --noconfirm go make gcc git
COPY . /app/
WORKDIR /app
RUN make build
ENTRYPOINT ["/app/pkgstats"]
