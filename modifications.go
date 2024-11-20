package main

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func ApplyModifications(resource *unstructured.Unstructured, modifications string) error {
	// Parse modifications
	mods := strings.Split(modifications, ",")
	for _, mod := range mods {
		kv := strings.SplitN(mod, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid modification '%s'", mod)
		}
		key, value := kv[0], kv[1]
		// Apply modifications based on key
		// For simplicity, let's handle labels and annotations
		switch {
		case strings.HasPrefix(key, "metadata.labels."):
			labelKey := strings.TrimPrefix(key, "metadata.labels.")
			labels, found, err := unstructured.NestedStringMap(resource.Object, "metadata", "labels")
			if err != nil {
				return err
			}
			if !found {
				labels = make(map[string]string)
			}
			labels[labelKey] = value
			unstructured.SetNestedStringMap(resource.Object, labels, "metadata", "labels")
		case strings.HasPrefix(key, "metadata.annotations."):
			annotationKey := strings.TrimPrefix(key, "metadata.annotations.")
			annotations, found, err := unstructured.NestedStringMap(resource.Object, "metadata", "annotations")
			if err != nil {
				return err
			}
			if !found {
				annotations = make(map[string]string)
			}
			annotations[annotationKey] = value
			unstructured.SetNestedStringMap(resource.Object, annotations, "metadata", "annotations")
		default:
			return fmt.Errorf("unsupported modification key '%s'", key)
		}
	}
	return nil
}
