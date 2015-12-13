package main


import (
	"log"
	"flag"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/fields"
	vClient "github.com/vulcand/vulcand/api"
	vPlugin "github.com/vulcand/vulcand/plugin"
	kClient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"
	"encoding/json"
	"time"
	"k8s.io/kubernetes/pkg/runtime"
	"github.com/vulcand/vulcand/engine"
	"strconv"
)



var (
	vServer string
	kServer string
	namespace string
	labelQuery string
)



func init() {
	flag.StringVar(&kServer, "kubernetes-api-server", "https://kubernetes.default", "Kubernetes server to proxy.")
	flag.StringVar(&vServer, "vulcand-server", "http://docker:8182", "Vulcand server for external loadbalancing.")
	flag.StringVar(&labelQuery, "query", "", "Label query for services to expose via the external loadbalancer.")
	flag.StringVar(&namespace, "namespace", "default", "Namespace in which to look for services to match")
}



func ensureEndpointConfiguredForVulcand(client *vClient.Client, endpoints api.Endpoints) {

	name := endpoints.Name + "." + endpoints.Namespace

	backend := engine.Backend{Id: name, Type: "http",
		Settings: engine.HTTPBackendSettings{
			Timeouts: engine.HTTPBackendTimeouts{Read: "5s", Dial: "5s", TLSHandshake:"10s"},
			KeepAlive: engine.HTTPBackendKeepAlive{Period: "30s", MaxIdleConnsPerHost: 12}}}

	err := client.UpsertBackend(backend)

	if err != nil {
		log.Println("Encountered error when creating backend in vulcand: " + err.Error())
	}


	// add servers to the backend
	for _, element := range endpoints.Subsets {
		for _, address := range element.Addresses {
			url := "http://" + address.IP + ":" + strconv.Itoa(element.Ports[0].Port)
			log.Println("Attempting to add server to backend " + backend.GetUniqueId().Id + " using url: " + url)
			err = client.UpsertServer(backend.GetUniqueId(), engine.Server{Id: url, URL: url}, time.Second * 5)
			if err != nil {
				log.Println("Encountered error when adding server to backend: " + err.Error())
			}
		}
	}

	// make sure there's a frontend
	frontend := engine.Frontend{Id: name, Type: "http", Route: "Path(`/`)", BackendId: name}
	err = client.UpsertFrontend(frontend, time.Second * 5)

	if err != nil {
		log.Println("Encountered error when creating frontend in vulcand: " + err.Error())
	}

}


func removeUnusedEndpointsFromVulcand(client *vClient.Client, endpoints api.Endpoints) {

	name := endpoints.Name + "." + endpoints.Namespace

	backend, err := client.GetBackend(engine.BackendKey{Id: name})

	if err != nil {
		log.Println("Encountered error when getting backend to remove unused endpoints: " + err.Error())
		return
	}

	servers, err := client.GetServers(engine.BackendKey{Id: name})

	if err != nil {
		log.Println("Encountered error when getting servers to remove unused endpoints: " + err.Error())
		return
	}

	js, _ := json.Marshal(backend)
	log.Println(js)

	js, _ = json.Marshal(servers)
	log.Println(js)
}


func deserialize(result runtime.Object) (api.Endpoints, error) {
	obj, _ := json.Marshal(result)
	var endpoints api.Endpoints
	err := json.Unmarshal(obj, &endpoints)
	if err != nil {
		log.Println("Could not unmarshal channel result into endpoint type: \n" + string(obj))
		return endpoints, err
	}
	return endpoints, nil
}

func main() {

	flag.Parse()
	log.Println("Connecting to kubernetes via url " + kServer)
	log.Println("Connecting to vulcand via url " + vServer)
	log.Println("Provided label query: " + labelQuery)
	log.Println("Observing endpoints within namespace: " + namespace)

	vClient := vClient.NewClient(vServer, vPlugin.NewRegistry())
	kClient, err := kClient.New(&kClient.Config{Host:kServer})

	if err != nil {
		log.Println("Error encountered when connecting to kubernetes api." + err.Error())
		panic(err)
	}

	var labelSelector labels.Selector = nil

	if labelQuery != "" {
		labelSelector, err = labels.Parse(labelQuery)
		if err != nil {
			log.Println("Error parsing the provided label query.")
			panic(err)
		}
	}  else {
		labelSelector = labels.Everything()
	}


	socket, err := kClient.Endpoints(namespace).
	Watch(labelSelector, fields.Everything(), api.ListOptions{Watch: true})

	if err != nil {
		log.Println("Error obtaining a watch on the kubernetes endpoints.")
		panic(err)
	}

	// poll the channel indefinitely
	for {

		select {
		case event := <-socket.ResultChan():
			switch event.Type {
			case watch.Added:
				endpoint, _ := deserialize(event.Object)
				ensureEndpointConfiguredForVulcand(vClient, endpoint)
				log.Println("Endpoint was added: \n" + endpoint.Name)
			case watch.Modified:
				endpoint, _ := deserialize(event.Object)
				ensureEndpointConfiguredForVulcand(vClient, endpoint)
				log.Println("Endpoint was modified: \n" + endpoint.Name)
			case watch.Deleted:
				endpoint, _ := deserialize(event.Object)
				removeUnusedEndpointsFromVulcand(vClient, endpoint)
				log.Println("Endpoint was deleted: \n" + endpoint.Name)
			case watch.Error:
				log.Println("Encountered an error from the endpoints socket. Continuing...")
			}
		default:
			time.Sleep(1 * time.Second)
		}

	}

}