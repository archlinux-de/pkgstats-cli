use anyhow::Error;
use clap::{crate_name, crate_version};
use serde::{de::DeserializeOwned, Serialize};
use std::sync::Arc;
use std::{borrow::Borrow, time::Duration};
use ureq::Agent;
use url::Url;

pub struct HttpClient {
    base_url: String,
    agent: Agent,
}

impl HttpClient {
    #[allow(clippy::result_large_err)]
    fn accept_json(
        req: ureq::Request,
        next: ureq::MiddlewareNext,
    ) -> Result<ureq::Response, ureq::Error> {
        next.handle(req.set("Accept", "application/json"))
    }

    fn create_client(base_url: &str, timeout: u64) -> Agent {
        ureq::AgentBuilder::new()
            .timeout(Duration::from_secs(timeout))
            .user_agent(&format!("{}/{}", crate_name!(), crate_version!()))
            .https_only(base_url.starts_with("https"))
            .middleware(Self::accept_json)
            .tls_connector(Arc::new(native_tls::TlsConnector::new().unwrap()))
            .build()
    }

    pub fn new(base_url: &str, timeout: u64) -> Self {
        Self {
            base_url: base_url.to_string(),
            agent: Self::create_client(base_url, timeout),
        }
    }

    pub fn query<T>(&self, path: &str) -> Result<T, Error>
    where
        T: DeserializeOwned,
    {
        self.agent
            .get(
                Url::parse([&self.base_url, path].concat().as_str())
                    .unwrap()
                    .as_str(),
            )
            .call()
            .map_err(Error::from)?
            .into_json()
            .map_err(Error::from)
    }

    pub fn query_with_params<T, I, K, V>(&self, path: &str, params: I) -> Result<T, Error>
    where
        T: DeserializeOwned,
        I: IntoIterator,
        I::Item: Borrow<(K, V)>,
        K: AsRef<str>,
        V: AsRef<str>,
    {
        self.agent
            .get(
                Url::parse_with_params([&self.base_url, path].concat().as_str(), params)
                    .unwrap()
                    .as_str(),
            )
            .call()
            .map_err(Error::from)?
            .into_json()
            .map_err(Error::from)
    }

    pub fn send_json<T>(&self, path: &str, payload: T) -> Result<String, Error>
    where
        T: Serialize,
    {
        self.agent
            .post(
                Url::parse([&self.base_url, path].concat().as_str())
                    .unwrap()
                    .as_str(),
            )
            .send_json(payload)
            .map_err(Error::from)?
            .into_string()
            .map_err(Error::from)
    }
}
