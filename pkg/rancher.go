package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type RancherConfig struct {
	URL   string
	Token string
}

func LoadRancherConfig() (*RancherConfig, error) {
	rancherURL := os.Getenv("RANCHER_URL")
	rancherToken := os.Getenv("RANCHER_TOKEN")

	if rancherURL == "" || rancherToken == "" {
		return nil, fmt.Errorf("RANCHER_URL and RANCHER_TOKEN environment variables are required")
	}

	return &RancherConfig{
		URL:   rancherURL,
		Token: rancherToken,
	}, nil
}

func GetClusterKubeconfig(config *RancherConfig, clusterID string) ([]byte, error) {
	url := fmt.Sprintf("%s/v3/clusters/%s?action=generateKubeconfig", config.URL, clusterID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the Bearer token
	req.Header.Set("Authorization", "Bearer "+config.Token)

	// Create a custom HTTP client that skips TLS verification
	insecureTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	insecureClient := &http.Client{Transport: insecureTransport}

	resp, err := insecureClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get kubeconfig: %s", string(body))
	}

	var result struct {
		Config string `json:"config"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return []byte(result.Config), nil
}
