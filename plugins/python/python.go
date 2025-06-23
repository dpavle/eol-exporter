//go:build plugin

package main

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type PythonPlugin struct{}

func (p *PythonPlugin) GetProductAndVersion() (product string, version string, err error) {

	cmd := exec.Command("python", "--version")

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	versionPattern := regexp.MustCompile(`Python [0-9]+.[0-9]+`)
	pythonVersion := strings.TrimPrefix(versionPattern.FindString(out.String()), "Python ")

	return "python", pythonVersion, err

}

var DataPlugin PythonPlugin
