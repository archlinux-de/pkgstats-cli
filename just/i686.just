test:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=386 go test -v -exec 'qemu-i386 -cpu coreduo' ./...
    # 32-Bit on x86_64
    CGO_ENABLED=0 GOARCH=386 go test -v -exec 'qemu-x86_64 /usr/bin/linux32' ./...

build:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=386 go build -buildvcs=false -o tests/build/pkgstats-i686

test-cpu-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=386 go run -exec 'qemu-i386 -cpu coreduo' main.go architecture system | grep -q '^i686$'
    # 32-Bit on x86_64
    CGO_ENABLED=0 GOARCH=386 go run -exec 'qemu-x86_64 /usr/bin/linux32' main.go architecture system | grep -q '^x86_64'

test-os-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=386 go run -exec 'qemu-i386 -cpu coreduo' main.go architecture os | grep -q '^i686$'
    # 32-Bit on x86_64
    CGO_ENABLED=0 GOARCH=386 go run -exec 'qemu-x86_64 /usr/bin/linux32' main.go architecture os | grep -q '^i686$'
