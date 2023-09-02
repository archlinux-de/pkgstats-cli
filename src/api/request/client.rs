use crate::api::http_client::HttpClient;
use anyhow::Error;
use rayon::prelude::*;
use serde::Deserialize;

pub struct Client {
    http_client: HttpClient,
}

#[derive(Debug, Deserialize)]
pub struct PackagePopularity {
    pub name: String,
    pub popularity: f32,
}

#[derive(Debug, Deserialize)]
pub struct PackagePopularityList {
    pub total: u32,
    pub count: u32,
    #[serde(rename = "packagePopularities")]
    pub package_popularities: Vec<PackagePopularity>,
}

impl Client {
    pub fn new(base_url: &str) -> Self {
        Self {
            http_client: HttpClient::new(base_url, 5),
        }
    }

    pub fn get_packages(&self, packages: &[String]) -> Result<PackagePopularityList, Error> {
        let mut package_popularities = packages
            .par_iter()
            .map(|package: &String| self.get_package(package))
            .collect::<Result<Vec<PackagePopularity>, Error>>()?;

        package_popularities.sort_unstable_by(|a, b| b.popularity.total_cmp(&a.popularity));

        Ok(PackagePopularityList {
            total: 0,
            count: 0,
            package_popularities,
        })
    }

    pub fn get_package(&self, package: &str) -> Result<PackagePopularity, Error> {
        self.http_client
            .query(&format!("/api/packages/{}", package))
    }

    pub fn search_packages(
        &self,
        query: &str,
        limit: &u16,
    ) -> Result<PackagePopularityList, Error> {
        self.http_client.query_with_params(
            "/api/packages",
            &[
                ("limit", limit.to_string()),
                ("offset", "0".to_string()),
                ("query", query.to_string()),
            ],
        )
    }
}
