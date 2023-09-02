use super::{Os, OsArchitecture};

impl OsArchitecture for Os {
    fn get_architecture() -> Option<String> {
        Some(std::env::consts::ARCH.to_string())
    }
}
