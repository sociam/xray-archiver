package main
import (
   "fmt"
   "flag"
   "log"
   "pipeline/util"
   "encoding/json"
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

var ip = flag.String("name", "www.cs.ox.ac.uk", "Name to lookup.")

func main() {

     flag.Parse()
     fmt.Printf("Tea config [%s] looking up %s \n\n", *cfgFile, *ip)
     blep()
     fmt.Printf("GeoIP Endpoint Configured [%s]\n", util.Cfg.GeoIPEndpoint);

     iptest := make([]string, 0, 0)
     for _, v := range util.MakeRange(0, 103) {
     	 iptest = append(iptest, *ip)
	 _ = v
     }

     result,err := util.GetHostGeoIPs(iptest)
     if (err != nil) {
     	fmt.Printf("fuckssake")
     }

     ms,merr := json.Marshal(result)
     if (merr != nil) { 
          fmt.Printf("merr is error: ",merr.Error())
     } else {
          fmt.Printf("result: ", string(ms))
     }

     for host, rs := range result {
     	 fmt.Println(" host ", host, " rs ", len(rs))
	 if len(rs) > 0 {
	    fmt.Println("rs[0].IP", rs[0].IP, rs[0].Latitude, rs[0].Longitude);
	 }
     }


}
