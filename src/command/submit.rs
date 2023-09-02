use crate::api::submit::{client::Client, request::Request};
use anyhow::Error;
use serde_json::json;

pub fn run(quiet: &bool, dump_json: &bool, base_url: &str) -> Result<(), Error> {
    if !*quiet && !*dump_json {
        println!("Collecting data...");
    }

    let request = Request::create();

    if *dump_json {
        println!("{}", json!(request));
    } else {
        if !*quiet {
            println!("Submitting data...");
        }

        let client = Client::new(base_url);
        client.send_request(&request)?;

        if !*quiet {
            println!("Data were successfully sent");
        }
    }

    Ok(())
}
