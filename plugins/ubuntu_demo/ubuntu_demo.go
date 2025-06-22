package main

type PrintMessagePlugin struct {}

func (p *PrintMessagePlugin) GetProductAndVersion() (product string, version string, err error) {

	product = "ubuntu"
	version = "22.04"
	err = nil
	return

}

var DataPlugin PrintMessagePlugin
