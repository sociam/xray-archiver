package main
import (
   "fmt"
   "flag"
   "log"
   "pipeline/util"
)

func blep() {
	var err error
	err = util.LoadCfg(*cfgFile, util.APIServ)
	if err != nil {
		log.Fatalf("Failed to read config: %s", err.Error())
	}
}

var cfgFile = flag.String("cfg", "/etc/xray/config.json", "config file location")
var port = flag.Uint("port", 8118, "Port to serve on.")
var ip = flag.String("name", "web.mit.edu", "Name to lookup.")

func main() {

     flag.Parse()
     fmt.Printf("Tea config [%s] looking up %s \n\n", *cfgFile, *ip)
     blep()
     fmt.Printf("GeoIP Endpoint Configured [%s]\n", util.Cfg.GeoIPEndpoint);

     result,err := util.GetHostGeoIP(*ip)
     if (err != nil) {
     	fmt.Printf("fuckssake")
     }
     fmt.Printf("Length [%d]\n", len(result))
     fmt.Printf("%f, %f\n", result[0].Latitude, result[0].Longitude)
}
