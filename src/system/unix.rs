use super::{Os, OsArchitecture};

impl OsArchitecture for Os {
    fn get_architecture() -> Option<String> {
        let machine = nix::sys::utsname::uname()
            .ok()?
            .machine()
            .to_string_lossy()
            .to_string();

        Some(machine)
    }
}
