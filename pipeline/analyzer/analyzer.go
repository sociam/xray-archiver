package main

import (
	"./config"
	"flag"
	"fmt"
	"net"
	"os"
	"path"
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

		split := strings.SplitN(s, "-", 4)
		if len(split) < 4 {
			fmt.Printf("failed to parse input \"%s\"\n", s)
		} else {
			app := App{split[0], split[1], split[2], split[3]}
			fmt.Printf("Got app %v\n", app)
			fmt.Print("Unpacking... ")
			err = unpack(app)
			if err != nil {
				fmt.Println()
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("done.")
			fmt.Print("Running simple analysis... ")
			_, err := simple_analyze(app)
			if err != nil {
				fmt.Println()
				fmt.Println(err.Error())
			}
			fmt.Println("done.")
		}
	}
}
