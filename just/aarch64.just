set quiet := true

export GOARCH := 'arm64'

[group('test')]
[working-directory('../')]
test:
    go test -exec qemu-aarch64 ./tests/...

[group('test')]
[working-directory('../')]
build:
    go build -buildvcs=false -o /dev/null

[group('test')]
[working-directory('../')]
test-cpu-detection:
    go run -exec 'qemu-aarch64' main.go architecture system | grep -q '^aarch64$'

[group('test')]
[working-directory('../')]
test-os-detection:
    go run -exec 'qemu-aarch64' main.go architecture os | grep -q '^aarch64$'
