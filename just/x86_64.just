test:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=amd64 go test -v -exec 'qemu-x86_64 -cpu Conroe' ./...
    CGO_ENABLED=0 GOARCH=amd64 go test -v -exec 'qemu-x86_64 -cpu Nehalem' ./...
    CGO_ENABLED=0 GOARCH=amd64 go test -v -exec 'qemu-x86_64 -cpu Haswell,-hle,-rtm' ./...

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
    CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64 -cpu Haswell,-hle,-rtm' main.go architecture system | grep -q '^x86_64_v3$'

test-os-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64 -cpu Conroe' main.go architecture os | grep -q '^x86_64$'
    CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64 -cpu Nehalem' main.go architecture os | grep -q '^x86_64$'
    CGO_ENABLED=0 GOARCH=amd64 go run -exec 'qemu-x86_64 -cpu Haswell,-hle,-rtm' main.go architecture os | grep -q '^x86_64$'
