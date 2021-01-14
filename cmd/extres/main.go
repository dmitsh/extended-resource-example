package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const resourceName = "example.com~1extres"

type patchExtendedResource struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value uint32 `json:"value"`
}

func main() {
	// get the node name
	nodeName := os.Getenv("NODE_NAME")
	log.Infof("NODE_NAME: %s", nodeName)
	// create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// get the node object
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err.Error())
	}
	// create and send patch request
	payload := []patchExtendedResource{{
		Op:    "add",
		Path:  "/status/capacity/" + resourceName,
		Value: getResourceCapacity(),
	}}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = clientset.CoreV1().Nodes().Patch(context.TODO(), node.Name, types.JSONPatchType, payloadBytes, metav1.PatchOptions{}, "status")
	if err != nil {
		log.Fatalln(err.Error())
	}

	var g run.Group
	{
		// Termination handler.
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		cancel := make(chan struct{})
		g.Add(
			func() error {
				select {
				case <-term:
					log.Infoln("Received SIGTERM, exiting gracefully...")
					onExit()
				case <-cancel:
					onExit()
				}
				return nil
			},
			func(err error) {
				close(cancel)
			},
		)
	}
	// add additional groups here
	if err := g.Run(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Infoln("Exit")
}

// getResourceCapacity return resource capacity on the node
func getResourceCapacity() uint32 {
	return 1073741824 // 1Gi
}

// onExit does necessary cleanup before exiting
func onExit() {

}
