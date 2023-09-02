use clap::{command, Parser, Subcommand};

const DEFAULT_BASE_URL: &str = "https://pkgstats.archlinux.de";
const SEARCH_DEFAULT_LIMIT: u16 = 10;
const SEARCH_MIN_LIMIT: i64 = 1;
const SEARCH_MAX_LIMIT: i64 = 10000;
const SHOW_MIN_PACKAGES: usize = 1;
const SHOW_MAX_PACKAGES: usize = 20;

#[derive(Parser)]
#[command(arg_required_else_help(true))]
pub struct Cli {
    #[arg(long, default_value = DEFAULT_BASE_URL, hide = true, help = "base url of the pkgstats server")]
    pub base_url: Option<String>,

    #[command(subcommand)]
    pub command: Option<Commands>,
}

#[derive(Debug, Subcommand)]
pub enum Commands {
    #[command(
        hide = true,
        alias = "arch",
        about = "Shows information about CPU and OS architecture"
    )]
    Architecture {
        #[command(subcommand)]
        command: ArchitectureCommands,
    },

    #[command(about = "Search packages and list their popularity")]
    Search {
        #[arg()]
        name: String,

        #[arg(short, long, default_value_t = SEARCH_DEFAULT_LIMIT, value_parser = clap::value_parser!(u16).range(SEARCH_MIN_LIMIT..SEARCH_MAX_LIMIT),help = format!("Limit the results from {} to {} entries", SEARCH_MIN_LIMIT, SEARCH_MAX_LIMIT))]
        limit: u16,
    },

    #[command(about = "Show one or more packages and compare their popularity")]
    Show {
        #[arg(num_args = SHOW_MIN_PACKAGES..SHOW_MAX_PACKAGES)]
        names: Vec<String>,
    },

    #[command(
        about = "Submit a list of your installed packages to the pkgstats project",
        long_about = format!("Submit a list of your installed packages, your system architecture\nand the mirror you are using to the pkgstats project.\n\nStatistics are available at {}", DEFAULT_BASE_URL)
    )]
    Submit {
        #[arg(
            short,
            long,
            default_value_t = false,
            exclusive = true,
            help = "Dump information that would be sent as JSON"
        )]
        dump_json: bool,

        #[arg(
            short,
            long,
            default_value_t = false,
            exclusive = true,
            help = "Suppress any output except errors"
        )]
        quiet: bool,
    },

    #[command(about = "Show the pkgstats client version")]
    Version {},
}

#[derive(Debug, Subcommand)]
pub enum ArchitectureCommands {
    #[command(about = "Shows OS architecture")]
    Os,

    #[command(alias = "cpu", about = "Shows CPU architecture")]
    System,
}

impl Cli {
    pub fn create() -> Self {
        Self::parse()
    }
}
