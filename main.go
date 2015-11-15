package main
import (
	"flag"
	"log"
	kClient "k8s.io/kubernetes/pkg/client/unversioned"
	vClient "github.com/vulcand/vulcand/api"
	vulcandPlugin "github.com/vulcand/vulcand/plugin"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/fields"
	"os"
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

	vClient := vClient.NewClient(vServer, vulcandPlugin.NewRegistry())
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

	results, error := kClient.Endpoints(namespace).List(labelSelector, fields.Everything())

	if error != nil {
		log.Println("Error parsing the provided label query.")
		os.Exit(1)
	}

	print(vClient)
	print(results)
	print(labelSelector)

	//	consumer, err := kClient.Endpoints(namespace).Watch(labelSelector, fields.Everything(), api.ListOptions{Watch: true})

	//	print("huh")
	//
	//	if err != nil {
	//		log.Println("Error encountered when getting a watch on the kubernetes endpoints resource.")
	//	}
	//
	//	print(consumer)

}