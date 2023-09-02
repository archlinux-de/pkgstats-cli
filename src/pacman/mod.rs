use std::process::Command;
use url::Url;

pub fn get_installed_packages() -> Vec<String> {
    let output = Command::new("pacman").arg("-Qq").output().unwrap();

    let s: String = String::from_utf8(output.stdout).unwrap();

    let r: Vec<String> = s.lines().map(str::to_string).collect();

    r
}

pub fn get_server() -> Option<String> {
    let output = Command::new("pacman-conf")
        .args(["--repo", "core", "Server"])
        .output()
        .ok()?;

    let server_string = String::from_utf8(output.stdout).ok()?;

    let mut server = Url::parse(&server_string).ok()?;
    // remove core/os/$arch from path
    server.path_segments_mut().unwrap().pop().pop().pop();

    server.set_query(None);
    server.set_fragment(None);

    server.set_username("").ok()?;
    server.set_password(None).ok()?;

    Some(server.to_string())
}
