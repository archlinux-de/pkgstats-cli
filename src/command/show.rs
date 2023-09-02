use crate::api::request::{
    client::Client,
    printer::{print_package_popularities, print_show_url},
    validator::validate_package_names,
};
use anyhow::{anyhow, Error};

pub fn run(names: &[String], base_url: &str) -> Result<(), Error> {
    if !validate_package_names(names) {
        return Err(anyhow!("invalid package names".to_string()));
    }

    let client = Client::new(base_url);

    let pacakge_popularities = client.get_packages(names)?;

    print_package_popularities(&pacakge_popularities);

    println!();

    let a: Vec<&str> = names.iter().map(String::as_str).collect();

    print_show_url(base_url, a.as_slice());

    Ok(())
}
