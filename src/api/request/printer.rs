use super::client::PackagePopularityList;
use percent_encoding::{utf8_percent_encode, NON_ALPHANUMERIC};
use std::io::Write;
use tabwriter::TabWriter;
use url::Url;

pub fn print_package_popularities(pacakge_popularities: &PackagePopularityList) {
    let mut tw = TabWriter::new(vec![]).minwidth(15).padding(0);

    for package_popularity in &pacakge_popularities.package_popularities {
        writeln!(
            &mut tw,
            "{}\t{:.2}",
            package_popularity.name, package_popularity.popularity
        )
        .unwrap();
    }
    tw.flush().unwrap();

    print!("{}", String::from_utf8(tw.into_inner().unwrap()).unwrap());

    if pacakge_popularities.count > 0 && pacakge_popularities.total > 0 {
        println!(
            "\n{} of {} results",
            pacakge_popularities.count, pacakge_popularities.total
        );
    }
}

pub fn print_search_url(base_url: &str, query: &str) {
    if !query.is_empty() {
        println!("See more results at {}/packages#query={}", base_url, query)
    }
}

pub fn print_show_url(base_url: &str, packages: &[&str]) {
    if !packages.is_empty() {
        let mut encoded_packages: Vec<String> = packages
            .iter()
            .map(|package| utf8_percent_encode(package, NON_ALPHANUMERIC).to_string())
            .collect();
        encoded_packages.sort();

        let mut link = Url::parse(base_url).unwrap();
        link.set_path("/compare/packages");
        link.set_fragment(Some(
            format!("packages={}", encoded_packages.join(",").as_str()).as_str(),
        ));

        println!("See more results at {}", link.as_str());
    }
}
