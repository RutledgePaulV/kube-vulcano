package main
import (
	"flag"
)



var (
	vulcandServer string
)

func init() {
	flag.StringVar(&vulcandServer, "vulcand-server", "http://docker:8182", "Vulcand server for external loadbalancing (ip:port)")
}

func main() {

	flag.Parse()
	print("hello")


}