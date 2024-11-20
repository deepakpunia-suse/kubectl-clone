package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Parse command-line flags
	var sourceCluster string
	var sourceNamespace string
	var sourceProject string
	var resourceType string
	var resourceName string
	var targetCluster string
	var targetNamespace string
	var targetProject string
	var newName string
	var modifications string

	flag.StringVar(&sourceCluster, "source-cluster", "", "Source cluster name or ID in Rancher")
	flag.StringVar(&sourceNamespace, "source-namespace", "", "Namespace of the resource to clone (default: default)")
	flag.StringVar(&sourceProject, "source-project", "", "Source project ID in Rancher")
	flag.StringVar(&resourceType, "type", "", "Type of the resource to clone (e.g., deployment, service)")
	flag.StringVar(&resourceName, "name", "", "Name of the resource to clone")
	flag.StringVar(&targetCluster, "target-cluster", "", "Target cluster name or ID in Rancher")
	flag.StringVar(&targetNamespace, "target-namespace", "", "Target namespace to clone the resource into")
	flag.StringVar(&targetProject, "target-project", "", "Target project ID in Rancher")
	flag.StringVar(&newName, "new-name", "", "New name for the cloned resource")
	flag.StringVar(&modifications, "modify", "", "Modifications to apply in key=value format, comma-separated")
	flag.Parse()

	// Validate required flags
	if sourceCluster == "" || resourceType == "" || resourceName == "" || targetCluster == "" {
		fmt.Println("Error: --source-cluster, --type, --name, and --target-cluster are required.")
		os.Exit(1)
	}

	if sourceNamespace == "" {
		sourceNamespace = "default"
	}

	if targetNamespace == "" {
		targetNamespace = sourceNamespace
	}

	// Load Rancher configuration
	rancherConfig, err := LoadRancherConfig()
	if err != nil {
		fmt.Printf("Error loading Rancher config: %v\n", err)
		os.Exit(1)
	}

	// Fetch kubeconfig for source and target clusters
	sourceKubeconfig, err := GetClusterKubeconfig(rancherConfig, sourceCluster)
	if err != nil {
		fmt.Printf("Error fetching source cluster kubeconfig: %v\n", err)
		os.Exit(1)
	}

	targetKubeconfig, err := GetClusterKubeconfig(rancherConfig, targetCluster)
	if err != nil {
		fmt.Printf("Error fetching target cluster kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create dynamic clients for source and target clusters
	sourceClient, err := GetDynamicClient(sourceKubeconfig)
	if err != nil {
		fmt.Printf("Error creating source cluster client: %v\n", err)
		os.Exit(1)
	}

	targetClient, err := GetDynamicClient(targetKubeconfig)
	if err != nil {
		fmt.Printf("Error creating target cluster client: %v\n", err)
		os.Exit(1)
	}

	// Get the GVR (GroupVersionResource) for the resource type
	gvr, err := GetGVR(resourceType)
	if err != nil {
		fmt.Printf("Error getting GVR: %v\n", err)
		os.Exit(1)
	}

	// Fetch the resource from the source cluster
	resource, err := sourceClient.Resource(gvr).Namespace(sourceNamespace).Get(Context(), resourceName, GetOptions())
	if err != nil {
		fmt.Printf("Error fetching resource from source cluster: %v\n", err)
		os.Exit(1)
	}

	// Modify the resource
	clonedResource := resource.DeepCopy()

	// Set new name
	if newName != "" {
		clonedResource.SetName(newName)
	} else {
		clonedResource.SetName(resourceName + "-clone")
	}

	// Set new namespace
	clonedResource.SetNamespace(targetNamespace)

	// Apply modifications
	if modifications != "" {
		err = ApplyModifications(clonedResource, modifications)
		if err != nil {
			fmt.Printf("Error applying modifications: %v\n", err)
			os.Exit(1)
		}
	}

	// Remove fields that should not be copied
	RemoveUnwantedFields(clonedResource)

	// Create the cloned resource in the target cluster
	_, err = targetClient.Resource(gvr).Namespace(targetNamespace).Create(Context(), clonedResource, CreateOptions())
	if err != nil {
		fmt.Printf("Error creating cloned resource in target cluster: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Resource '%s' cloned successfully as '%s' in cluster '%s', namespace '%s'.\n", resourceName, clonedResource.GetName(), targetCluster, targetNamespace)
}
