mod api;
mod cli;
mod command;
mod pacman;
mod system;

use anyhow::Error;
use cli::Commands;
use command::{architecture, search, show, submit, version};

/*
 * @TODO
 * Split up commands into seperate files
 *     https://stackoverflow.com/questions/73357989/clap-subcommands-over-multiple-rs-files
 *     https://blog.cloudnativefolks.org/writing-rust-clis-clap
 *     https://blog.railway.app/p/rust-cli-rewrite
 *     https://github.com/sile-typesetter/casile/blob/master/Cargo.toml
 * Possible ureq replacement: https://github.com/sbstp/attohttpc
 * Disable timer and submit for unsupported arches
 * Support riscv
 * Use rust cross to cross compile and test
 * Implement arch without subcommands
 * Use git-describe as version
 * See https://github.com/rust-lang/stdarch/blob/master/crates/std_detect/src/detect/os/linux/arm.rs
 * See https://crates.io/crates/cupid
 *
 * Test tokio and https://hyper.rs/guides/1/client/basic/
 *      https://github.com/hyperium/hyper/blob/master/examples/client_json.rs
 */

const DEFAULT_BASE_URL: &str = "https://pkgstats.archlinux.de";

fn main() -> Result<(), Error> {
    let cli = cli::Cli::create();

    let mut base_url: String = DEFAULT_BASE_URL.to_string();
    if let Some(url) = cli.base_url.as_deref() {
        base_url = url.to_string();
    }

    match &cli.command {
        Some(Commands::Architecture { command }) => architecture::run(command),
        Some(Commands::Search { name, limit }) => search::run(name, limit, &base_url)?,
        Some(Commands::Show { names }) => show::run(names, &base_url)?,
        Some(Commands::Submit { dump_json, quiet }) => submit::run(quiet, dump_json, &base_url)?,
        Some(Commands::Version {}) => version::run(),
        None => unreachable!(),
    }

    Ok(())
}
