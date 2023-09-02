use crate::{
    cli::ArchitectureCommands,
    system::{get_architecture, get_cpu_architecture},
};

pub fn run(command: &ArchitectureCommands) {
    match command {
        ArchitectureCommands::Os => println!("{}", get_architecture().unwrap()),
        ArchitectureCommands::System => println!("{}", get_cpu_architecture().unwrap()),
    }
}
