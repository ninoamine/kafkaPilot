package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
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

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Failed to create dynamic client: %v", err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "kafkapilot.io",
		Version:  "v1alpha1",
		Resource: "topics",
	}

	namespace := "default"

	watchCRD(dynamicClient, gvr, namespace)

}

func watchCRD(client dynamic.Interface, gvr schema.GroupVersionResource, namespace string) {
	watcher, err := client.Resource(gvr).Namespace(namespace).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Fatalf("Failed to set up watch on CRD: %v", err)
	}

	fmt.Println("Watching for CRD events...")
	for event := range watcher.ResultChan() {
		obj, ok := event.Object.(runtime.Object)
		if !ok {
			klog.Warning("Unexpected type")
			continue
		}

		switch event.Type {
		case watch.Added:
			fmt.Println("CRD Added:", obj)
		case watch.Modified:
			fmt.Println("CRD Modified:", obj)
		case watch.Deleted:
			fmt.Println("CRD Deleted:", obj)
		default:
			fmt.Println("Unknown event type:", event.Type)
		}
	}
}

func getConfig(kubeconfig string) (*rest.Config, error) {
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


