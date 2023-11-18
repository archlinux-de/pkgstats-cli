use crate::{
    pacman,
    system::{get_architecture, get_cpu_architecture},
};
use serde::Serialize;

#[derive(Debug, Serialize)]
pub struct System {
    architecture: String,
}

#[derive(Debug, Serialize)]
pub struct OS {
    architecture: String,
}

#[derive(Debug, Serialize)]
pub struct Pacman {
    mirror: Option<String>,
    packages: Vec<String>,
}

#[derive(Debug, Serialize)]
pub struct Request {
    version: String,
    system: System,
    os: OS,
    pacman: Pacman,
}

impl Request {
    pub fn create() -> Self {
        Self {
            version: "3".to_owned(),
            system: System {
                architecture: get_cpu_architecture().unwrap(),
            },
            os: OS {
                architecture: get_architecture().unwrap(),
            },
            pacman: Pacman {
                mirror: pacman::get_server(),
                packages: pacman::get_installed_packages(),
            },
        }
    }
}
