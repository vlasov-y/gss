# Go Static Server

[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/9757/badge)](https://www.bestpractices.dev/projects/9757)
![Build](https://github.com/vlasov-y/gss/workflows/Build/badge.svg?branch=main)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

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
