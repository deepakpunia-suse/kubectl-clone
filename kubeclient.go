package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func GetDynamicClient(kubeconfigBytes []byte) (dynamic.Interface, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

func Context() context.Context {
	return context.TODO()
}

func GetOptions() metav1.GetOptions {
	return metav1.GetOptions{}
}

func CreateOptions() metav1.CreateOptions {
	return metav1.CreateOptions{}
}

func RemoveUnwantedFields(resource *unstructured.Unstructured) {
	unstructured.RemoveNestedField(resource.Object, "metadata", "resourceVersion")
	unstructured.RemoveNestedField(resource.Object, "metadata", "uid")
	unstructured.RemoveNestedField(resource.Object, "metadata", "creationTimestamp")
	unstructured.RemoveNestedField(resource.Object, "status")
}
