package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))
	for _, node := range nodes.Items {
		fmt.Printf("Found node %q %#v\n", node.Name, node)

		payload := []patchExtendedResource{{
			//Op:    "add",
			Op:    "replace",
			Path:  "/status/capacity/example.com~1mydev",
			Value: 10,
		}}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Patching with %q\n", string(payloadBytes))
		//result := clientset.CoreV1().RESTClient().Patch(types.JSONPatchType).Resource("nodes").Name(framework.TestContext.NodeName).SubResource("status").Body(patch).Do(context.TODO())
		n, err := clientset.CoreV1().Nodes().Patch(context.TODO(), node.Name, types.JSONPatchType, payloadBytes, metav1.PatchOptions{}, "status")
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Patched %#v\n", *n)
	}

	for {
		time.Sleep(10 * time.Second)
	}
}
