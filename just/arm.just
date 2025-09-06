set quiet := true

export GOARCH := 'arm'

[group('test')]
[working-directory('../')]
test:
    GOARM=5 go test -exec 'qemu-arm -cpu pxa250' ./tests/...
    GOARM=6 go test -exec 'qemu-arm -cpu arm1176' ./tests/...
    GOARM=7 go test -exec 'qemu-arm -cpu cortex-a15' ./tests/...
    go test -exec 'qemu-arm -cpu max' ./tests/...

[group('test')]
[working-directory('../')]
build:
    go build -buildvcs=false -o /dev/null

[group('test')]
[working-directory('../')]
test-cpu-detection:
    GOARM=5 go run -exec 'qemu-arm -cpu pxa250' main.go architecture system | grep -q '^armv5$'
    GOARM=6 go run -exec 'qemu-arm -cpu arm1176' main.go architecture system | grep -q '^armv6$'
    GOARM=7 go run -exec 'qemu-arm -cpu cortex-a15' main.go architecture system | grep -q '^armv7$'
    go run -exec 'qemu-arm -cpu max' main.go architecture system | grep -q '^aarch64$'

[group('test')]
[working-directory('../')]
test-os-detection:
    GOARM=5 go run -exec 'qemu-arm -cpu pxa250' main.go architecture os | grep -q '^armv5tel$'
    GOARM=6 go run -exec 'qemu-arm -cpu arm1176' main.go architecture os | grep -q '^armv6l$'
    GOARM=7 go run -exec 'qemu-arm -cpu cortex-a15' main.go architecture os | grep -q '^armv7l$'
    go run -exec 'qemu-arm -cpu max' main.go architecture os | grep -q '^armv7l$'
