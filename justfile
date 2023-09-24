# list all recipes
default:
	@just --list

# Prepare sources in order to build offline
prepare:
	cargo fetch --locked

# build pkgstats for production
build:
	cargo build --frozen --release

# update dependencies
update:
	cargo update

# run cargo clippy
check-clippy:
	cargo clippy

# check cargo fmt
check-fmt:
	cargo fmt --check

# run all static checks
check: check-fmt check-clippy

# run unit tests
test:
	cargo --frozen test

# See https://github.com/cross-rs/cross
# and https://github.com/cargo-bins/cargo-binstall
#     https://github.com/marketplace/actions/install-development-tools
#
# TODO: use cargo install?
#
# separate justfile for cross testing
# use test martix on ci

cargo-aarch64 command *options:
	cross {{command}} --target aarch64-unknown-linux-musl -F bundled-tls {{options}}

cargo-armv7 command *options:
	cross {{command}} --target armv7-unknown-linux-musleabihf -F bundled-tls {{options}}

cargo-i686 command *options:
	cross {{command}} --target i686-unknown-linux-musl -F bundled-tls {{options}}

cargo-riscv64 command *options:
	# musl libc is currently unavailable on riscv64
	cross {{command}} --target riscv64gc-unknown-linux-gnu -F bundled-tls {{options}}

cargo-x86_64 command *options:
	cross {{command}} --target x86_64-unknown-linux-musl -F bundled-tls {{options}}

# run unit tests on different CPU architectures
test-cross-platform: (cargo-aarch64 'test') (cargo-armv7 'test') (cargo-i686 'test') (cargo-riscv64 'test') (cargo-x86_64 'test')

# build for different CPU architectures
test-build: (cargo-aarch64 'build') (cargo-armv7 'build') (cargo-i686 'build') (cargo-riscv64 'build') (cargo-x86_64 'build')

test-cpu-detection-armv7:
	just cargo-armv7 run -- architecture system | grep -q '^arm$'

test-cpu-detection-aarch64:
	just cargo-aarch64 run -- architecture system | grep -q '^aarch64$'

test-cpu-detection-riscv64:
	just cargo-riscv64 run -- architecture system | grep -q '^riscv64$'

test-cpu-detection-i686: (cargo-i686 'build')
	qemu-i386 -cpu coreduo target/i686-unknown-linux-musl/debug/pkgstats architecture system | grep -q '^i686$'
	@# i686 on x86_64
	qemu-x86_64 /usr/bin/linux32 target/i686-unknown-linux-musl/debug/pkgstats architecture system | grep -q '^x86_64'

test-cpu-detection-x86_64: (cargo-x86_64 'build')
	@# x86_64
	qemu-x86_64 -cpu Conroe target/x86_64-unknown-linux-musl/debug/pkgstats architecture system | grep -q '^x86_64$'
	qemu-x86_64 -cpu Nehalem target/x86_64-unknown-linux-musl/debug/pkgstats architecture system | grep -q '^x86_64_v2$'
	@# Test crashes on older Qemu versions
	if qemu-x86_64 -version | grep -Eq 'version (7\.[2-9]|[8-9]\.)[0-9]*\.[0-9]+$'; then qemu-x86_64 -cpu Haswell target/x86_64-unknown-linux-musl/debug/pkgstats architecture system 2>&1 | grep -q '^x86_64_v3$'; fi

# test cpu architecture detection on different CPUs
test-cpu-detection: test-cpu-detection-armv7 test-cpu-detection-aarch64 test-cpu-detection-riscv64 test-cpu-detection-x86_64 test-cpu-detection-i686

test-os-detection-aarch64:
	just cargo-aarch64 run -- architecture os | grep -q '^aarch64$'

test-os-detection-armv7:
	just cargo-armv7 run -- architecture os | grep -q '^armv7l$'

test-os-detection-i686: (cargo-i686 'build')
	qemu-i386 -cpu coreduo target/i686-unknown-linux-musl/debug/pkgstats architecture os | grep -q '^i686$'
	@# i686 on x86_64
	qemu-x86_64 /usr/bin/linux32 target/i686-unknown-linux-musl/debug/pkgstats architecture os | grep -q '^i686$'

test-os-detection-riscv64:
	just cargo-riscv64 run -- architecture os | grep -q '^riscv64$'

test-os-detection-x86_64:
	just cargo-x86_64 run -- architecture os | grep -q '^x86_64$'

# test os architecture detection on different CPUs
test-os-detection: test-os-detection-aarch64 test-os-detection-armv7 test-os-detection-i686 test-os-detection-riscv64 test-os-detection-x86_64

# run integration tests with a mocked API server
test-integration:
	docker buildx build --pull . -f tests/integration/Dockerfile -t pkgstats-test-integration

# install pkgstats and its configuration
install *DESTDIR='':
	@# cli
	install -D target/release/pkgstats -m755 "{{DESTDIR}}/usr/bin/pkgstats"

	#@TODO: install only on supported targets

	@# systemd timer
	for service in pkgstats.service pkgstats.timer; do \
		install -Dt "{{DESTDIR}}/usr/lib/systemd/system" -m644 init/${service} ; \
	done
	install -d "{{DESTDIR}}/usr/lib/systemd/system/timers.target.wants"
	cd "{{DESTDIR}}/usr/lib/systemd/system/timers.target.wants" && ln -s ../pkgstats.timer

	@# bash completions
	install -D target/release/completions/pkgstats.bash "{{DESTDIR}}/usr/share/bash-completion/completions/pkgstats"

	@# zsh completions
	install -Dt "{{DESTDIR}}/usr/share/zsh/site-functions/" target/release/completions/_pkgstats

	@# fish completions
	install -Dt "{{DESTDIR}}/usr/share/fish/vendor_completions.d" target/release/completions/pkgstats.fish

# run all available tests
test-all: check test test-build test-cpu-detection test-os-detection test-integration

# remove any untracked and generated files
clean:
	git clean -fdqx -e .idea -e .vscode

# vim: set ft=make :
