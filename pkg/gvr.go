package main

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGVR(resourceType string) (schema.GroupVersionResource, error) {
	resourceType = strings.ToLower(resourceType)
	switch resourceType {
	case "deployment", "deployments":
		return schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}, nil
	case "service", "services":
		return schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}, nil
	case "configmap", "configmaps":
		return schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}, nil
	case "secret", "secrets":
		return schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}, nil
	// Add more resource types as needed
	default:
		return schema.GroupVersionResource{}, fmt.Errorf("unsupported resource type '%s'", resourceType)
	}
}
