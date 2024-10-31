test:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=386 GO386=softfloat go test -exec 'qemu-i386 -cpu pentium2' ./tests/...
    CGO_ENABLED=0 GOARCH=386 go test -exec 'qemu-i386 -cpu coreduo' ./tests/...
    # 32-Bit on x86_64
    CGO_ENABLED=0 GOARCH=386 go test -exec 'qemu-x86_64 /usr/bin/linux32' ./tests/...

build:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=386 go build -buildvcs=false -o /dev/null
    CGO_ENABLED=0 GOARCH=386 GO386=softfloat go build -buildvcs=false -o /dev/null

test-cpu-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=386 GO386=softfloat go run -exec 'qemu-i386 -cpu pentium' main.go architecture system | grep -q '^i586$'
    CGO_ENABLED=0 GOARCH=386 GO386=softfloat go run -exec 'qemu-i386 -cpu pentium2' main.go architecture system | grep -q '^i686$'
    CGO_ENABLED=0 GOARCH=386 go run -exec 'qemu-i386 -cpu coreduo' main.go architecture system | grep -q '^i686$'
    # 32-Bit on x86_64
    CGO_ENABLED=0 GOARCH=386 go run -exec 'qemu-x86_64 /usr/bin/linux32' main.go architecture system | grep -q '^x86_64'

test-os-detection:
    #!/usr/bin/env bash
    set -euo pipefail
    cd ../
    CGO_ENABLED=0 GOARCH=386 GO386=softfloat go run -exec 'qemu-i386 -cpu pentium' main.go architecture os | grep -q '^i586$'
    CGO_ENABLED=0 GOARCH=386 GO386=softfloat go run -exec 'qemu-i386 -cpu pentium2' main.go architecture os | grep -q '^i686$'
    CGO_ENABLED=0 GOARCH=386 go run -exec 'qemu-i386 -cpu coreduo' main.go architecture os | grep -q '^i686$'
    # 32-Bit on x86_64
    CGO_ENABLED=0 GOARCH=386 go run -exec 'qemu-x86_64 /usr/bin/linux32' main.go architecture os | grep -q '^i686$'
