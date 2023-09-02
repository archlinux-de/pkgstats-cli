use crate::api::request::{
    client::Client,
    printer::{print_package_popularities, print_search_url},
    validator::validate_package_name,
};
use anyhow::{anyhow, Error};

pub fn run(name: &str, limit: &u16, base_url: &str) -> Result<(), Error> {
    if !validate_package_name(name) {
        return Err(anyhow!("invalid package name"));
    }

    let client = Client::new(base_url);

    let pacakge_popularities = client.search_packages(name, limit)?;

    print_package_popularities(&pacakge_popularities);

    println!();

    print_search_url(base_url, name);

    Ok(())
}
