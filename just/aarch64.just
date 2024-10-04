test:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=arm64 go test -v -exec qemu-aarch64 ./...

build:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=arm64 go build -buildvcs=false -o tests/build/pkgstats-aarch64

test-cpu-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=arm64 go run -exec 'qemu-aarch64' main.go architecture system | grep -q '^aarch64$'

test-os-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=arm64 go run -exec 'qemu-aarch64' main.go architecture os | grep -q '^aarch64$'
