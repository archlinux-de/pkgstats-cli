use super::request::Request;
use crate::api::http_client::HttpClient;
use anyhow::Error;

pub struct Client {
    http_client: HttpClient,
}

impl Client {
    pub fn new(base_url: &str) -> Self {
        Self {
            http_client: HttpClient::new(base_url, 10),
        }
    }

    pub fn send_request(&self, request: &Request) -> Result<String, Error> {
        self.http_client.send_json("/api/submit", request)
    }
}
