test:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=arm GOARM=5 go test -v -exec 'qemu-arm -cpu pxa250' ./...
    CGO_ENABLED=0 GOARCH=arm GOARM=6 go test -v -exec 'qemu-arm -cpu arm1176' ./...
    CGO_ENABLED=0 GOARCH=arm GOARM=7 go test -v -exec 'qemu-arm -cpu cortex-a15' ./...
    CGO_ENABLED=0 GOARCH=arm go test -v -exec 'qemu-arm -cpu max' ./...

build:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=arm go build -buildvcs=false -a -o /dev/null

test-cpu-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=arm GOARM=5 go run -exec 'qemu-arm -cpu pxa250' main.go architecture system | grep -q '^armv5$'
    CGO_ENABLED=0 GOARCH=arm GOARM=6 go run -exec 'qemu-arm -cpu arm1176' main.go architecture system | grep -q '^armv6$'
    CGO_ENABLED=0 GOARCH=arm GOARM=7 go run -exec 'qemu-arm -cpu cortex-a15' main.go architecture system | grep -q '^armv7$'
    CGO_ENABLED=0 GOARCH=arm go run -exec 'qemu-arm -cpu max' main.go architecture system | grep -q '^aarch64$'

test-os-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=arm GOARM=5 go run -exec 'qemu-arm -cpu pxa250' main.go architecture os | grep -q '^armv5tel$'
    CGO_ENABLED=0 GOARCH=arm GOARM=6 go run -exec 'qemu-arm -cpu arm1176' main.go architecture os | grep -q '^armv6l$'
    CGO_ENABLED=0 GOARCH=arm GOARM=7 go run -exec 'qemu-arm -cpu cortex-a15' main.go architecture os | grep -q '^armv7l$'
    CGO_ENABLED=0 GOARCH=arm go run -exec 'qemu-arm -cpu max' main.go architecture os | grep -q '^armv7l$'
