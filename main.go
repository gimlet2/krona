package main

import (
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/ericchiang/k8s"
	"context"
	"log"
	"io/ioutil"
	"fmt"
	"github.com/ghodss/yaml"
	"strconv"
	"github.com/robfig/cron"
	"time"
	"net/http"
	"os"
)

type CronJobTriggerResource struct {
	Metadata *metav1.ObjectMeta `json:"metadata"`
	Spec CronSpec `json:"spec"`
}

type CronSpec struct {
	FunctionName string `json:"function-name"`
    Schedule string `json:"schedule"`
}

type FunctionResource struct {
	Metadata *metav1.ObjectMeta `json:"metadata"`
	Spec     Spec               `json:"spec"`
}

type Spec struct {
	Type     string  `json:"type"`
	Service  Service `json:"service"`
}
type Service struct {
	Ports []Port `json:"ports"`
}

type Port struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

func (m *FunctionResource) GetMetadata() *metav1.ObjectMeta {
	return m.Metadata
}

func (m *CronJobTriggerResource) GetMetadata() *metav1.ObjectMeta {
	return m.Metadata
}

type FunctionResourceList struct {
	Metadata *metav1.ListMeta   `json:"metadata"`
	Items    []FunctionResource `json:"items"`
}

type CronJobTriggerResourceList struct {
	Metadata *metav1.ListMeta   `json:"metadata"`
	Items    []CronJobTriggerResource `json:"items"`
}

// Require for MyResourceList to implement k8s.ResourceList
func (m *FunctionResourceList) GetMetadata() *metav1.ListMeta {
	return m.Metadata
}

func (m *CronJobTriggerResourceList) GetMetadata() *metav1.ListMeta {
	return m.Metadata
}

func loadClient(kubeconfigPath string) (*k8s.Client, error) {
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

type CronJob struct {
	Pattern string
	Cron    *cron.Cron
}

func main() {
	var crons = map[string]*CronJob{}
	log.Print("Start")
	// kubeConfigPath := "/home/user/.kube/config"
	kubeConfigPath := os.Getenv("KUBE_CONFIG") //"/home/user/.kube/config"
	var client *k8s.Client
	var err error
	if kubeConfigPath != "" {
		client, err = loadClient(kubeConfigPath)
	} else {
		client, err = k8s.NewInClusterClient()
	}
	if err != nil {
		log.Fatal(err)
	}
	k8s.Register("kubeless.io", "v1beta1", "functions", true, &FunctionResource{})
	k8s.Register("kubeless.io", "v1beta1", "cronjobtriggers", true, &CronJobTriggerResource{})
	k8s.RegisterList("kubeless.io", "v1beta1", "functions", true, &FunctionResourceList{})
	k8s.RegisterList("kubeless.io", "v1beta1", "cronjobtriggers", true, &CronJobTriggerResourceList{})
	
	var cronTriggers CronJobTriggerResourceList
	for {
		err = client.List(context.Background(), k8s.AllNamespaces, &cronTriggers)
		for f := range cronTriggers.Items {
			cronTrigger := cronTriggers.Items[f]
			var function FunctionResource
			err = client.Get(context.Background(), *cronTrigger.Metadata.Namespace, cronTrigger.Spec.FunctionName, &function)
			if err == nil {
				url := "http://" + *function.Metadata.Name + "." + *function.Metadata.Namespace + ":" + strconv.Itoa(function.Spec.Service.Ports[0].Port)
				
				pattern := cronTrigger.Spec.Schedule
				if crons[*function.Metadata.Name] != nil && crons[*function.Metadata.Name].Pattern != pattern {
					crons[*function.Metadata.Name].Cron.Stop()
					delete(crons, *function.Metadata.Name)
				}
				if crons[*function.Metadata.Name] == nil || crons[*function.Metadata.Name].Pattern != pattern {
					log.Printf("Function '%s' with schedule '%s' was descovered", *function.Metadata.Name, pattern)
					c := cron.New()
					c.AddFunc(pattern, func() {
						log.Printf("Trigger '%s' function - GET - %s", *function.Metadata.Name, url)
						http.Get(url)
					})
					crons[*function.Metadata.Name] = &CronJob{pattern, c}
					c.Start()
				}
			} else {
				log.Print(err);
			}
			
		}
		time.Sleep(5 * time.Second)
	}
}
