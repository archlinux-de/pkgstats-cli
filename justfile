set quiet := true

import 'just/dev.just'

export CGO_CPPFLAGS := env_var_or_default('CPPFLAGS', '')
export CGO_CFLAGS := env_var_or_default('CFLAGS', '')
export CGO_CXXFLAGS := env_var_or_default('CXXFLAGS', '')
export CGO_LDFLAGS := env_var_or_default('LDFLAGS', '')
export CGO_ENABLED := '0'

# list all recipes
[private]
default:
    @just --list

# Prepare sources in order to build offline
[group('install')]
prepare:
    go mod download

# build pkgstats for production
[group('install')]
build:
    CGO_ENABLED=1 go build -a -o pkgstats \
       	-buildmode=pie -mod=readonly -modcacherw -buildvcs=false \
       	-ldflags '-compressdwarf=false -linkmode=external -s -w -X pkgstats-cli/internal/build.Version={{ `git describe --tags` }}'

# run unit tests
[group('test')]
test:
    go test ./tests/...

# install pkgstats and its configuration
[group('install')]
install *DESTDIR='':
    #!/usr/bin/env bash
    set -euo pipefail

    # cli
    install -D pkgstats -m755 "{{ DESTDIR }}/usr/bin/pkgstats"

    # systemd timer
    for service in pkgstats.{service,timer}; do
       	install -Dt "{{ DESTDIR }}/usr/lib/systemd/system" -m644 init/${service}
    done
    install -d "{{ DESTDIR }}/usr/lib/systemd/system/timers.target.wants"
    ln -s ../pkgstats.timer -t "{{ DESTDIR }}/usr/lib/systemd/system/timers.target.wants"

    # bash completions
    install -d "{{ DESTDIR }}/usr/share/bash-completion/completions"
    ./pkgstats completion bash > "{{ DESTDIR }}/usr/share/bash-completion/completions/pkgstats"

    # zsh completions
    install -d "{{ DESTDIR }}/usr/share/zsh/site-functions/"
    ./pkgstats completion zsh > "{{ DESTDIR }}/usr/share/zsh/site-functions/_pkgstats"

    # fish completions
    install -d "{{ DESTDIR }}/usr/share/fish/vendor_completions.d"
    ./pkgstats completion fish > "{{ DESTDIR }}/usr/share/fish/vendor_completions.d/pkgstats.fish"
