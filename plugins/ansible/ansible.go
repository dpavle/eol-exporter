//go:build plugin

package main

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type AnsiblePlugin struct{}

func (p *AnsiblePlugin) GetProductAndVersion() (product string, version string, err error) {

	cmd := exec.Command("ansible", "--version")

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	versionPattern := regexp.MustCompile(`core [0-9]+.[0-9]+`)
	ansibleVersion := strings.TrimPrefix(versionPattern.FindString(out.String()), "core ")

	return "ansible-core", ansibleVersion, err

}

var DataPlugin AnsiblePlugin
