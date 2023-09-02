trait OsArchitecture {
    fn get_architecture() -> Option<String>;
}

struct Os;
struct Sys;

trait SystemArchitecture {
    fn get_cpu_architecture() -> Option<String>;
}

#[cfg(target_family = "unix")]
#[path = "unix.rs"]
mod os;

#[cfg(not(target_family = "unix"))]
#[path = "os.rs"]
mod os;

pub fn get_architecture() -> Option<String> {
    Os::get_architecture()
}

#[cfg(any(
    all(target_arch = "x86", target_feature = "sse"),
    target_arch = "x86_64"
))]
#[path = "x86.rs"]
mod sys;

#[cfg(not(any(
    all(target_arch = "x86", target_feature = "sse"),
    target_arch = "x86_64"
)))]
#[path = "sys.rs"]
mod sys;

pub fn get_cpu_architecture() -> Option<String> {
    Sys::get_cpu_architecture()
}

#[cfg(test)]
mod tests {
    #[cfg(all(target_arch = "x86_64"))]
    #[test]
    fn it_gets_cpu_architecture() {
        use crate::system::get_cpu_architecture;
        assert!(get_cpu_architecture().unwrap().starts_with("x86_64"));
    }

    #[cfg(all(target_arch = "x86_64"))]
    #[test]
    fn it_gets_architecture() {
        use crate::system::get_architecture;
        assert_eq!("x86_64", get_architecture().unwrap());
    }
}
