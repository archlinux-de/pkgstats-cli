use super::{Sys, SystemArchitecture};

impl SystemArchitecture for Sys {
    fn get_cpu_architecture() -> Option<String> {
        Some(std::env::consts::ARCH.to_string())
    }
}
