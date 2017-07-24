package main

import (
	"./config"
	"flag"
	"fmt"
)

var cfgFile = flag.String("cfg", "/etc/xray/config.json", "config file location")
var cfg config.Config

func init() {
	flag.Parse()
	cfg = config.Load(*cfgFile)
}

func main() {
	fmt.Println("Starting xray analyzer daemon")
	runServer()
}
