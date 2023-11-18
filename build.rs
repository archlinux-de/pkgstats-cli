use clap::crate_name;
use clap::CommandFactory;
use clap_complete::{generate_to, Shell};
use std::{env, fs, path};

include!("src/cli.rs");

fn main() {
    set_version();
    generate_shell_completion();
}

fn generate_shell_completion() {
    let outdir = match env::var_os("OUT_DIR") {
        None => return,
        Some(outdir) => outdir,
    };

    if path::Path::new(&outdir).ancestors().count() < 3 {
        return;
    }

    let completions_dir = path::Path::new(&outdir)
        .parent()
        .unwrap()
        .parent()
        .unwrap()
        .parent()
        .unwrap()
        .join("completions");

    fs::create_dir_all(&completions_dir).unwrap();

    let mut cmd = Cli::command();
    for &shell in &[Shell::Bash, Shell::Fish, Shell::Zsh] {
        generate_to(shell, &mut cmd, crate_name!(), &completions_dir).unwrap();
    }
}

fn set_version() {
    if let Ok(val) = std::env::var("VERSION") {
        println!("cargo:rustc-env=CARGO_PKG_VERSION={}", val);
    }
    println!("cargo:rerun-if-env-changed=VERSION");
}
