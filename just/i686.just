set quiet := true

export GOARCH := '386'

[group('test')]
[working-directory('../')]
test:
    GO386=softfloat go test -exec 'qemu-i386 -cpu pentium2' ./tests/...
    go test -exec 'qemu-i386 -cpu coreduo' ./tests/...
    # 32-Bit on x86_64
    go test -exec 'qemu-x86_64 /usr/bin/linux32' ./tests/...

[group('test')]
[working-directory('../')]
build:
    go build -buildvcs=false -o /dev/null
    GO386=softfloat go build -buildvcs=false -o /dev/null

[group('test')]
[working-directory('../')]
test-cpu-detection:
    GO386=softfloat go run -exec 'qemu-i386 -cpu pentium' main.go architecture system | grep -q '^i586$'
    GO386=softfloat go run -exec 'qemu-i386 -cpu pentium2' main.go architecture system | grep -q '^i686$'
    go run -exec 'qemu-i386 -cpu coreduo' main.go architecture system | grep -q '^i686$'
    # 32-Bit on x86_64
    go run -exec 'qemu-x86_64 /usr/bin/linux32' main.go architecture system | grep -q '^x86_64'

[group('test')]
[working-directory('../')]
test-os-detection:
    GO386=softfloat go run -exec 'qemu-i386 -cpu pentium' main.go architecture os | grep -q '^i586$'
    GO386=softfloat go run -exec 'qemu-i386 -cpu pentium2' main.go architecture os | grep -q '^i686$'
    go run -exec 'qemu-i386 -cpu coreduo' main.go architecture os | grep -q '^i686$'
    # 32-Bit on x86_64
    go run -exec 'qemu-x86_64 /usr/bin/linux32' main.go architecture os | grep -q '^i686$'
