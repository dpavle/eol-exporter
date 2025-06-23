//go:build plugin

package main

import (
	"os/exec"
	"bytes"
	"log"
	"regexp"
)

type DockerPlugin struct {}

func (p *DockerPlugin) GetProductAndVersion() (product string, version string, err error) {

	cmd := exec.Command("docker", "version", "--format", "{{ .Server.Version }}")

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	versionPattern := regexp.MustCompile(`^[0-9]+.[0-9]+`)
	dockerVersion := versionPattern.FindString(out.String())

	return "docker-engine", dockerVersion, err

}

var DataPlugin DockerPlugin
