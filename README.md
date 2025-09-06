# pkgstats

`pkgstats` is the command-line client for the [Arch Linux package statistics project](https://pkgstats.archlinux.de/). It allows users to submit a list of their installed packages, helping Arch Linux developers understand package usage and prioritize their efforts.

The tool also provides functionality to search for packages and compare their popularity based on the collected data.

## For End-Users

### Installation

You can install `pkgstats` from the official Arch Linux repositories:

```bash
sudo pacman -Syu pkgstats
```

### Usage

#### Submitting Data

The package will collect the list of your installed packages, your system's architecture, and the mirror you are using, and submit it to the pkgstats project. The data is sent anonymously.

You can also view the data that would be sent without actually submitting it:

```bash
pkgstats submit --dump-json
```

#### Searching for Packages

You can search for a package to see its popularity:

```bash
pkgstats search <package-name>
```

Example:

```bash
pkgstats search firefox
```

#### Comparing Package Popularity

You can compare the popularity of multiple packages:

```bash
pkgstats show <package1> <package2> ...
```

Example:

```bash
pkgstats show firefox chromium
```

## For Developers

This project uses `just` as a command runner. To see all available commands, run:

```bash
just
```

### Building from Source

To build `pkgstats` from source, you need to have `go` and `just` installed.

1.  Clone the repository:

    ```bash
    git clone https://github.com/pkgstats/pkgstats-cli.git
    cd pkgstats-cli
    ```

2.  Build the project:

    ```bash
    just build
    ```

    This will create a `pkgstats` binary in the root directory.

### Testing

The project has a comprehensive test suite that includes unit tests, integration tests, and static code analysis.

#### Static Analysis

To run all static analysis checks, including formatting, vetting, and linting, run:

```bash
just lint
```

#### Unit Tests

To run the unit tests, run:

```bash
just test
```

This will run all tests in the `tests/` directory.

To generate a test coverage report, run:

```bash
just coverage
```

#### Cross-platform Tests

The project includes tests for different CPU architectures. To run them, you need to have `docker` and `qemu` installed.

- `just test-cross-platform`: Runs unit tests on different CPU architectures.
- `just test-build`: Builds the project for different CPU architectures.
- `just test-cpu-detection`: Tests CPU architecture detection on different CPUs.
- `just test-os-detection`: Tests OS architecture detection on different CPUs.

#### Integration Tests

To run the integration tests, you need to have `docker` installed. The integration tests run with a mocked API server.

```bash
just test-integration
```

#### All Tests

To run all available tests, including static analysis, unit tests, and integration tests, run:

```bash
just test-all
```

### Justfile Setup

The `justfile` setup is modular to support different CPU architectures.

- The main `justfile` imports `just/dev.just`, which contains the main development tasks.
- `just/dev.just` includes `just` files for each supported architecture (`aarch64`, `arm`, `i686`, `loongarch64`, `riscv64`, `x86_64`) using the `mod` keyword.
- Each architecture-specific `just` file defines how to run tests and builds for that architecture using `qemu`.
- The cross-platform testing tasks in `just/dev.just` iterate through all supported architectures and execute the corresponding tasks from the architecture-specific `just` files.

### Architecture

The `pkgstats` project is structured as follows:

- `cmd/`: Contains the command-line interface logic, using the `cobra` library. Each command has its own file.
- `internal/`: Contains the core logic of the application.
  - `api/`: Handles communication with the pkgstats API.
  - `pacman/`: Interacts with the pacman configuration to gather package information.
  - `system/`: Gathers system information like CPU architecture.
- `main.go`: The main entry point of the application.
- `justfile`: Contains the `just` commands for development and testing.
- `tests/`: Contains the unit and integration tests.
