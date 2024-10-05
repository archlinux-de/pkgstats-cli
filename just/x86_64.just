test:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=amd64 go test -v -exec 'qemu-x86_64' ./...

build:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=amd64 go build -buildvcs=false -a -o /dev/null

test-cpu-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64 -cpu Conroe' main.go architecture system | grep -q '^x86_64$'
    CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64 -cpu Nehalem' main.go architecture system | grep -q '^x86_64_v2$'
    # Test crashes on older Qemu versions
    if qemu-x86_64 -version | grep -Eq 'version (7\.[2-9]|[8-9]\.)[0-9]*\.[0-9]+$'; then CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64 -cpu Haswell' main.go architecture system 2>&1 | grep -q '^x86_64_v3$'; fi

test-os-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64' main.go architecture os | grep -q '^x86_64$'
