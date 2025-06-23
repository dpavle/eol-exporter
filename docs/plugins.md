# Plugins for eol-exporter

The eol-exporter supports a Go native plugin system, allowing you to monitor End-of-Life (EOL) status for additional software and products beyond the OS and kernel. This document explains how to configure, use, and develop plugins for eol-exporter, using the existing plugins as reference.

---

## What is a Plugin?

A plugin in eol-exporter is a Go shared object (`.so`) library that provides product and version detection for installed software. Plugins are loaded at runtime and must implement a standard interface. The exporter uses this information to fetch EOL data from [endoflife.date](https://endoflife.date/).

Plugins must be placed under the `plugins/` directory, in a subdirectory named after the plugin. For example, the `ansible` plugin should be located at `plugins/ansible/ansible.so`.

---

## Using Plugins

### 1. Configuration

To use plugins, list their names in your `eol-exporter` configuration file (specified with `--config`).  
Example configuration section:

```yaml
plugins:
  - ansible
  - docker
  - python
```

Each name refers to a plugin shared object expected at `plugins/{name}/{name}.so`.

---

## Writing a Plugin

Plugins must be written in Go and compiled as shared objects. Each plugin must provide a type with a `GetProductAndVersion()` method and export a variable named `DataPlugin` of that type.

### Minimal Plugin Example

```go
//go:build plugin

package main

import (
	"os/exec"
	"bytes"
	"log"
	"regexp"
	"strings"
)

type ExamplePlugin struct{}

func (p *ExamplePlugin) GetProductAndVersion() (product string, version string, err error) {
	cmd := exec.Command("examplecmd", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	versionPattern := regexp.MustCompile(`[0-9]+\\.[0-9]+`)
	exampleVersion := versionPattern.FindString(out.String())
	return "example", exampleVersion, err
}

var DataPlugin ExamplePlugin
```

#### Reference Plugins

- [Ansible plugin example](https://github.com/dpavle/eol-exporter/blob/main/plugins/ansible/ansible.go)
- [Docker plugin example](https://github.com/dpavle/eol-exporter/blob/main/plugins/docker/docker.go)
- [Python plugin example](https://github.com/dpavle/eol-exporter/blob/main/plugins/python/python.go)

### Compiling a Plugin

```bash
go build -buildmode=plugin -tags plugin -o plugins/your_plugin/your_plugin.so path/to/your_plugin.go
```
The plugin must be named `{plugin}.so` and placed in `plugins/{plugin}/`.

---

## Plugin Discovery

- The exporter loads each plugin specified in the configuration from `plugins/{name}/{name}.so`.
- If a plugin fails to load, it will be skipped and a warning will be logged.

---

## Plugin Lifecycle

- Plugins are loaded at application startup and invoked at each scrape (every time Prometheus scrapes the exporter).
- Plugins should be fast and lightweight.
- Execution timeout: 10 seconds (subject to change/configuration).

---

## Troubleshooting

- Ensure `.so` files are built for the target system's architecture and Go version.
- Verify plugin paths and names match those listed in the configuration.
- Review eol-exporter logs for errors related to plugin loading or execution.

---

## Best Practices

- Use simple and unique names for your plugins.
- Avoid plugins that require high privileges or may impact system stability.
- Keep plugin logic simple and focused on version detection.
- Document any plugin-specific configuration or prerequisites.

---

## See Also

- [endoflife.date products list](https://endoflife.date/api/products.json)
- [eol-exporter README](../README.md)
- Go [plugin documentation](https://pkg.go.dev/plugin)
- [All plugins in eol-exporter repository](https://github.com/dpavle/eol-exporter/tree/main/plugins)

---

_If you have questions or want to share your plugin, please open an issue or pull request!_

---

**Note:**  
This list of reference plugins may be incomplete. For more examples and the latest plugins, [browse plugins/ in the GitHub repository](https://github.com/dpavle/eol-exporter/tree/main/plugins).
