package main

import (
	"./config"
	"flag"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

var cfgFile = flag.String("cfg", "/etc/xray/config.json", "config file location")
var cfg config.Config

func init() {
	flag.Parse()
	cfg = config.Load(*cfgFile)
}

func checkDir(dir, name string) {
	fif, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0644)
			if err != nil {
				//TODO: something else
				panic(fmt.Sprintf("Couldn't create %s: %s", name, err.Error()))
			}
		} else {
			//TODO: something else
			panic(err)
		}
	}

	if !fif.IsDir() {
		panic(fmt.Sprintf("%s isn't a directory!", name))
	}
}

func main() {
	fmt.Println("Starting xray analyzer daemon")

	checkDir(cfg.UnpackDir, "Unpacked APK directory")
	checkDir(path.Dir(cfg.SockPath), "Socket directory")

	_ = os.Remove(cfg.SockPath) // probably catch errors here?

	apkSockAddr, err := net.ResolveUnixAddr("unixgram", cfg.SockPath)
	if err != nil {
		//TODO: not this
		panic(err)
	}
	apkSock, err := net.ListenUnixgram("unixgram", apkSockAddr)
	if err != nil {
		//TODO: not this
		panic(err)
	}
	defer apkSock.Close()

	for {
		b := make([]byte, 1024)
		n, err := apkSock.Read(b)
		if err != nil {
			//TODO: not this
			panic(err)
		}
		s := string(b[:n])

		if s == "" {
			// maybe do something else
			os.Exit(0)
		}

		split := strings.SplitN(s, "-", 5)
		if len(split) < 5 {
			fmt.Printf("failed to parse input \"%s\"\n", s)
		} else {
			id, err := strconv.Atoi(split[0])
			if err != nil {
				fmt.Printf("Given app id %s isn't a number!\n", id)
				continue
			}
			app := App{id, split[1], split[2], split[3], split[4]}
			fmt.Printf("Got app %v\n", app)
			fmt.Print("Unpacking... ")
			err = unpack(app)
			if err != nil {
				fmt.Println()
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("done.")
			fmt.Println("Running simple analysis... ")
			out, err := simple_analyze(app)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Printf("Hosts found: %v\n\n", out)
		}
	}
}
