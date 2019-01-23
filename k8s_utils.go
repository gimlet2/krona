package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ericchiang/k8s"
	"github.com/ghodss/yaml"
)

// Client is wrapper on top of k8s.Client
type Client interface {
	GetCronTriggers() CronJobTriggerResourceList
	GetFunction(namespace string, name string) (FunctionResource, error)
}

// ClientImpl implementation of Client interface
type ClientImpl struct {
	*k8s.Client
}

// NewClient constructor of Client
func NewClient(kubeConfigPath string) Client {
	var client *k8s.Client
	var err error
	if kubeConfigPath != "" {
		client, err = loadClientWithConfig(kubeConfigPath)
	} else {
		client, err = k8s.NewInClusterClient()
	}
	if err != nil {
		log.Fatal(err)
	}
	return &ClientImpl{client}
}

func loadClientWithConfig(kubeconfigPath string) (*k8s.Client, error) {
	data, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("read kubeconfig: %v", err)
	}

	// Unmarshal YAML into a Kubernetes config object.
	var config k8s.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal kubeconfig: %v", err)
	}
	return k8s.NewClient(&config)
}

// GetCronTriggers fetchs CronTriggers
func (c ClientImpl) GetCronTriggers() CronJobTriggerResourceList {
	var cronTriggers CronJobTriggerResourceList
	err := c.List(context.Background(), k8s.AllNamespaces, &cronTriggers)
	if err != nil {
		log.Printf("Failed to fetch cronTriggers - %v", err)
	}
	return cronTriggers
}

// GetFunction fetchs Function by name
func (c ClientImpl) GetFunction(namespace string, name string) (FunctionResource, error) {
	var function FunctionResource
	err := c.Get(context.Background(), namespace, name, &function)
	if err != nil {
		log.Printf("Failed to fetch function - %v", err)
		return function, err
	}
	return function, nil
}

// RegisterResources registers custom resources
func RegisterResources() {
	k8s.Register("kubeless.io", "v1beta1", "functions", true, &FunctionResource{})
	k8s.Register("kubeless.io", "v1beta1", "cronjobtriggers", true, &CronJobTriggerResource{})
	k8s.RegisterList("kubeless.io", "v1beta1", "functions", true, &FunctionResourceList{})
	k8s.RegisterList("kubeless.io", "v1beta1", "cronjobtriggers", true, &CronJobTriggerResourceList{})
}
