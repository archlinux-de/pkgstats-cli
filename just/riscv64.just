set quiet := true

export GOARCH := 'riscv64'

[group('test')]
[working-directory('../')]
test:
    go test -exec qemu-riscv64 ./tests/...

[group('test')]
[working-directory('../')]
build:
    go build -buildvcs=false -o /dev/null

[group('test')]
[working-directory('../')]
test-cpu-detection:
    go run -exec 'qemu-riscv64' main.go architecture system | grep -q '^riscv64$'

[group('test')]
[working-directory('../')]
test-os-detection:
    go run -exec 'qemu-riscv64' main.go architecture os | grep -q '^riscv64$'
