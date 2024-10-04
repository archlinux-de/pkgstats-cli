test:
    #!/usr/bin/env bash
    set -euo pipefail
    command -v qemu-loongarch64 > /dev/null 2>&1 || exit 0
    cd ../
    CGO_ENABLED=0 GOARCH=loong64 go test -v -exec qemu-loongarch64 ./...

build:
    #!/usr/bin/env bash
    set -euo pipefail
    command -v qemu-loongarch64 > /dev/null 2>&1 || exit 0
    cd ../
    CGO_ENABLED=0 GOARCH=loong64 go build -buildvcs=false -o tests/build/pkgstats-loongarch64

test-cpu-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    command -v qemu-loongarch64 > /dev/null 2>&1 || exit 0
    cd ../
    CGO_ENABLED=0 GOARCH=loong64 go run -exec 'qemu-loongarch64' main.go architecture system | grep -q '^loong64$'

test-os-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    command -v qemu-loongarch64 > /dev/null 2>&1 || exit 0
    cd ../
    CGO_ENABLED=0 GOARCH=loong64 go run -exec 'qemu-loongarch64' main.go architecture os | grep -q '^loongarch64$'
