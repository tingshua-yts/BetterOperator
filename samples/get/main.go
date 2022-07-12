package main

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := "/Users/tingshuai.yts/.kube/config.online"
	var (
		config *rest.Config
		err    error
	)

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating client: %v", err)
		os.Exit(1)
	}

	client := kubernetes.NewForConfigOrDie(config)

	list, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing nodes: %v", err)
		os.Exit(1)
	}

	for _, node := range list.Items {
		fmt.Printf("Node: %s\n", node.Name)
	}
}
