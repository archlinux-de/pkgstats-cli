pub fn validate_package_name(name: &str) -> bool {
    if name.is_empty() || name.len() > 190 {
        return false;
    }

    // see https://gitlab.archlinux.org/pacman/pacman/-/blob/master/scripts/libmakepkg/lint_pkgbuild/pkgname.sh.in#L32
    name.starts_with(|c| matches!(c, '0'..='9' | 'A'..='Z' | 'a'..='z' | '@' | '_' | '+'))
        && name[1..]
            .chars()
            .all(|c| matches!(c, '0'..='9' | 'A'..='Z' | 'a'..='z' | '@' | '.' | '_' | '+' | '-'))
}

pub fn validate_package_names(names: &[String]) -> bool {
    names.iter().all(|name| validate_package_name(name))
}

#[cfg(test)]
mod tests {
    use crate::api::request::validator::{validate_package_name, validate_package_names};

    #[test]
    fn it_validates_package_name() {
        const CASES: &[(&str, bool)] = &[
            ("pacman", true),
            ("@pacman", true),
            ("_pacman", true),
            ("+pacman", true),
            ("pacman-contrib", true),
            ("pacman_foo", true),
            ("pacman@7", true),
            ("pacman+bar", true),
            ("pacman.bar", true),
            ("-pacman", false),
            (".pacman", false),
            ("ö", false),
            ("", false),
        ];

        for &(name, valid) in CASES.iter() {
            assert_eq!(valid, validate_package_name(&name));
        }
    }

    #[test]
    fn it_validates_package_names() {
        let cases: &[(&[String], bool)] = &[
            (&["pacman".to_string(), "pacman-contrib".to_string()], true),
            (&["pacman".to_string(), "-pacman".to_string()], false),
            (&["-pacman".to_string()], false),
            (&["pacman".to_string()], true),
        ];

        for &(names, valid) in cases.iter() {
            assert_eq!(valid, validate_package_names(names));
        }
    }
}
