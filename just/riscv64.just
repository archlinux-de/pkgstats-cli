test:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=riscv64 go test -v -exec qemu-riscv64 ./...

build:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=riscv64 go build -buildvcs=false -o tests/build/pkgstats-riscv64

test-cpu-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=riscv64 go run -exec 'qemu-riscv64' main.go architecture system | grep -q '^riscv64$'

test-os-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=riscv64 go run -exec 'qemu-riscv64' main.go architecture os | grep -q '^riscv64$'
