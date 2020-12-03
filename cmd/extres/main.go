package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type patchExtendedResource struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value uint32 `json:"value"`
}

func main() {
	// get the node name
	nodeName := os.Getenv("NODE_NAME")
	fmt.Printf("NodeName: %s\n", nodeName)
	// create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	// get the node object
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Found node %q %s\n", node.Name, node.String())
	// create and send patch request
	payload := []patchExtendedResource{{
		Op:    "add",
		Path:  "/status/capacity/example.com~1mydev",
		Value: 1073741824, // 1 Gi
	}}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err.Error())
	}
	//result := clientset.CoreV1().RESTClient().Patch(types.JSONPatchType).Resource("nodes").Name(framework.TestContext.NodeName).SubResource("status").Body(patch).Do(context.TODO())
	_, err = clientset.CoreV1().Nodes().Patch(context.TODO(), node.Name, types.JSONPatchType, payloadBytes, metav1.PatchOptions{}, "status")
	if err != nil {
		panic(err.Error())
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
					fmt.Println("Received SIGTERM, exiting gracefully...")
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
		panic(err.Error())
	}
	fmt.Println("Exit")
}

func onExit() {
	// on-exit housekeeping
}
