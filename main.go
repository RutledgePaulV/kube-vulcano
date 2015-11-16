package main


import (
	"os"
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



func main() {

	flag.Parse()
	log.Println("Connecting to kubernetes via url " + kServer)
	log.Println("Connecting to vulcand via url " + vServer)
	log.Println("Provided label query: " + labelQuery)
	log.Println("Observing endpoints within namespace: " + namespace)

	vClient := vClient.NewClient(vServer, vPlugin.NewRegistry())
	kClient, err := kClient.NewInCluster()

	if err != nil {
		log.Println("Error encountered when connecting to kubernetes api.")
		os.Exit(1)
	}

	var labelSelector labels.Selector = nil

	if labelQuery != "" {
		labelSelector, err = labels.Parse(labelQuery)
		if err != nil {
			log.Println("Error parsing the provided label query.")
			os.Exit(1)
		}
	}  else {
		labelSelector = labels.Everything()
	}

	consumer, err := kClient.Endpoints(namespace).Watch(labelSelector, fields.Everything(), api.ListOptions{Watch: true})

	if err != nil {
		log.Println("Error obtaining a watch on the kubernetes endpoints.")
		os.Exit(1)
	}

	print(vClient)

	channel := consumer.ResultChan()

	for {
		time.Sleep(3 * time.Second)
		next := <- channel
		switch next.Type {
		case watch.Added:
			obj, _ := json.Marshal(next.Object)
			log.Println("Endpoint was added. \n" + string(obj))
			break
		case watch.Modified:
			obj, _ := json.Marshal(next.Object)
			log.Println("Endpoint was added. \n" + string(obj))
			break
		case watch.Deleted:
			obj, _ := json.Marshal(next.Object)
			log.Println("Endpoint was added. \n" + string(obj))
			break
		case watch.Error:
			obj, _ := json.Marshal(next.Object)
			log.Println("Endpoint was added. \n" + string(obj))
			break
		}
	}


}