package main
import (
	"flag"
	"github.com/vulcand/vulcand/api"
	"github.com/vulcand/vulcand/plugin"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"time"
	"encoding/json"
)

type Object struct {
	Object Endpoints `json:"object"`
}

type Endpoints struct {
	Kind       string   `json:"kind"`
	ApiVersion string   `json:"apiVersion"`
	Metadata   Metadata `json:"metadata"`
	Subsets    []Subset `json:"subsets"`
}

type Metadata struct {
	Name string `json:"name"`
}

type Subset struct {
	Addresses []Address `json:"addresses"`
	Ports     []Port    `json:"ports"`
}

type Address struct {
	IP string `json:"ip"`
}

type Port struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}



var (
	apiServer string
	vulcandServer string
)

func init() {
	flag.StringVar(&apiServer, "kube-apiserver", "127.0.0.1:8001", "Kubernetes API server for watching endpoints. (ip:port)")
	flag.StringVar(&vulcandServer, "vulcand-server", "192.168.99.100:8182", "Vulcand server for external loadbalancing (ip:port)")
}


func main() {

	flag.Parse()


	var client = api.NewClient(fmt.Sprintf("http://%s", vulcandServer), plugin.NewRegistry())


	origin := "http://localhost"
	url := fmt.Sprintf("ws://%s/api/v1/endpoints?watch=true", apiServer)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}


	for {
		var ep Object
		if err := websocket.JSON.Receive(ws, &ep); err != nil {
			log.Println(err)
			time.Sleep(time.Duration(2 * time.Second))
		}

		// get vulcand backends
		nr, err := client.GetBackends()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(2 * time.Second))
		}

		var js, _ = json.Marshal(nr)

		print(string(js))

	}

}