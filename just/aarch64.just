test:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=arm64 go test -exec qemu-aarch64 ./tests/...

build:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=arm64 go build -buildvcs=false -o /dev/null

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
