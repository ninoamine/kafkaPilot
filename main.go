package main

import (
	"flag"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func main() {

	klog.InitFlags(nil)
	defer klog.Flush()

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", home+"/.kube/config", "")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "")
	}

	flag.Parse()

	config, err := getConfig(*kubeconfig)
	if err != nil {
		klog.Fatalf("Failed to get kubeconfig: %v", err)
	}

}

func getConfig(kubeconfig string) (*rest.config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
