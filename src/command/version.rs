use clap::{crate_name, crate_version};

pub fn run() {
    println!("{}, version {}", crate_name!(), crate_version!());
}
