# Go Static Server

[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/9757/badge)](https://www.bestpractices.dev/projects/9757)
![Build](https://github.com/vlasov-y/gss/workflows/Build/badge.svg)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Coverage Status](https://badge.coveralls.io/repos/github/vlasov-y/gss/badge.svg?branch=main)](https://badge.coveralls.io/github/vlasov-y/gss?branch=main)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=vlasov-y_gss&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=vlasov-y_gss)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=vlasov-y_gss&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=vlasov-y_gss)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=vlasov-y_gss&metric=coverage)](https://sonarcloud.io/summary/new_code?id=vlasov-y_gss)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=vlasov-y_gss&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=vlasov-y_gss)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=vlasov-y_gss&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=vlasov-y_gss)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=vlasov-y_gss&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=vlasov-y_gss)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=vlasov-y_gss&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=vlasov-y_gss)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=vlasov-y_gss&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=vlasov-y_gss)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=vlasov-y_gss&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=vlasov-y_gss)

Simple server written in Go to serve a static site.

## Features

1. HTTP/HTTPS support
2. Client TLS auth
3. TLS curves, cipher suites and min/max version configuration
4. GZip compression
5. Custom headers (CORS supported)
6. Issuing certificate with ACME
7. Prometheus metrics
8. Configuration over both YAML/JSON and environment variables
