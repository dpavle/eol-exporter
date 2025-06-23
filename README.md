# eol-exporter

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dpavle/eol-exporter)](https://golang.org/doc/devel/release.html)
[![Build Status](https://github.com/dpavle/eol-exporter/actions/workflows/go.yml/badge.svg)](https://github.com/dpavle/eol-exporter/actions)
[![Latest Release](https://img.shields.io/github/v/release/dpavle/eol-exporter.svg?label=release)](https://github.com/dpavle/eol-exporter/releases)
[![Prometheus Exporter](https://img.shields.io/badge/prometheus-exporter-orange.svg)](https://prometheus.io/)

**eol-exporter** is a Prometheus exporter that exposes End-of-Life (EOL) information about your operating system, kernel, and installed software. By default, it collects and exports EOL data for your OS and kernel.

Additional software or products can be included via a simple plugin system, allowing you to monitor EOL status for any software or product.

The EOL data is pulled from the [endoflife.date API](https://endoflife.date/docs/api/v1/) and is refreshed every 24 hours.

> **Partially inspired by** [reimlima/endoflife_exporter](https://github.com/reimlima/endoflife_exporter)

---

## Installation

Build from source:

```bash
git clone https://github.com/dpavle/eol-exporter.git
cd eol-exporter
go build -o eol-exporter
```

Or download a pre-built binary from [Releases](https://github.com/dpavle/eol-exporter/releases).

## Usage

Start the exporter with default settings:

```bash
./eol-exporter
```

Options:

- `--config` — Path to config file
- `--listen-port` — Port to start HTTP exporter on (default: 3020)
- `--listen-address` — Address to start HTTP exporter on (default: 0.0.0.0)

Example:

```bash
./eol-exporter --listen-port=3020 --listen-address=127.0.0.1
```

## Prometheus Scrape Config

Add to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'eol-exporter'
    static_configs:
      - targets: ['localhost:3020']
```

## Plugins

Extend monitoring to additional software/products by adding plugins. See the [plugins documentation](docs/plugins.md) for details.

## Example Metrics

```
# HELP eol_date End of life date for the product release cycle. Expressed in seconds since Unix epoch (Unix Timestamp).
# TYPE eol_date gauge
eol_date{codename="",host="fedora",label="6.14",name="6.14",product="linux"} 1.7495136e+09
eol_date{codename="Adams",host="fedora",label="42 (Adams)",name="42",product="fedora"} 1.7786304e+09
# HELP product_details_info Full details of a product.
# TYPE product_details_info gauge
product_details_info{category="os",host="fedora",label="Fedora Linux",name="fedora",product="fedora",versionCommand="cat /etc/fedora-release"} 1
product_details_info{category="os",host="fedora",label="Linux Kernel",name="linux",product="linux",versionCommand="uname -r"} 1
# HELP product_release_info Full information about a product release cycle.
# TYPE product_release_info gauge
product_release_info{codename="",host="fedora",isEoas="false",isEoes="false",isEol="true",isLts="false",isMaintained="false",label="6.14",latest_link="https://kernelnewbies.org/Linux_6.14",latest_name="6.14.11",name="6.14",product="linux"} 1
product_release_info{codename="Adams",host="fedora",isEoas="false",isEoes="false",isEol="false",isLts="false",isMaintained="true",label="42 (Adams)",latest_link="",latest_name="",name="42",product="fedora"} 1
# HELP release_date Release date of the product release cycle. Expressed in seconds since Unix epoch (Unix Timestamp).
# TYPE release_date gauge
release_date{codename="",host="fedora",label="6.14",name="6.14",product="linux"} 1.7427744e+09
release_date{codename="Adams",host="fedora",label="42 (Adams)",name="42",product="fedora"} 1.7446752e+09
```

## Data Source

- [endoflife.date API](https://endoflife.date/docs/api/v1/)

## License

This project is licensed under the GNU GPLv3 License. See [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome!

---

_Acknowledgements: [endoflife.date](https://endoflife.date/) for lifecycle data._
