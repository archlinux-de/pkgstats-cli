set quiet := true

export GOARCH := 'amd64'

[group('test')]
[working-directory('../')]
test:
    go test -exec 'qemu-x86_64 -cpu Conroe' ./tests/...
    go test -exec 'qemu-x86_64 -cpu Nehalem' ./tests/...
    go test -exec 'qemu-x86_64 -cpu Haswell,-hle,-rtm' ./tests/...

[group('test')]
[working-directory('../')]
build:
    go build -buildvcs=false -o /dev/null

[group('test')]
[working-directory('../')]
test-cpu-detection:
    go run -exec 'qemu-x86_64 -cpu Conroe' main.go architecture system | grep -q '^x86_64$'
    go run -exec 'qemu-x86_64 -cpu Nehalem' main.go architecture system | grep -q '^x86_64_v2$'
    go run -exec 'qemu-x86_64 -cpu Haswell,-hle,-rtm' main.go architecture system | grep -q '^x86_64_v3$'

[group('test')]
[working-directory('../')]
test-os-detection:
    go run -exec 'qemu-x86_64 -cpu Conroe' main.go architecture os | grep -q '^x86_64$'
    go run -exec 'qemu-x86_64 -cpu Nehalem' main.go architecture os | grep -q '^x86_64$'
    go run -exec 'qemu-x86_64 -cpu Haswell,-hle,-rtm' main.go architecture os | grep -q '^x86_64$'
