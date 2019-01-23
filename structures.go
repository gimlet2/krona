package main

import (
	"strconv"

	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"github.com/robfig/cron"
)

// CronJobTriggerResource structure
type CronJobTriggerResource struct {
	Metadata *metav1.ObjectMeta `json:"metadata"`
	Spec     CronSpec           `json:"spec"`
}

// CronSpec structure
type CronSpec struct {
	FunctionName string `json:"function-name"`
	Schedule     string `json:"schedule"`
}

// FunctionResource structure
type FunctionResource struct {
	Metadata *metav1.ObjectMeta `json:"metadata"`
	Spec     Spec               `json:"spec"`
}

// Spec structure
type Spec struct {
	Type    string  `json:"type"`
	Service Service `json:"service"`
}

// Service structure
type Service struct {
	Ports []Port `json:"ports"`
}

// Port structure
type Port struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

// FunctionResourceList structure
type FunctionResourceList struct {
	Metadata *metav1.ListMeta   `json:"metadata"`
	Items    []FunctionResource `json:"items"`
}

// CronJobTriggerResourceList structure
type CronJobTriggerResourceList struct {
	Metadata *metav1.ListMeta         `json:"metadata"`
	Items    []CronJobTriggerResource `json:"items"`
}

// CronJob structure
type CronJob struct {
	Pattern string
	Cron    *cron.Cron
}

// GetMetadata function to get metadata
func (f *FunctionResource) GetMetadata() *metav1.ObjectMeta {
	return f.Metadata
}

// GetMetadata function to get metadata
func (m *CronJobTriggerResource) GetMetadata() *metav1.ObjectMeta {
	return m.Metadata
}

// GetPattern function to get job's pattern
func (m *CronJobTriggerResource) GetPattern() string {
	return m.Spec.Schedule
}

// GetFunctionName function to get job's function name
func (m *CronJobTriggerResource) GetFunctionName() string {
	return m.Spec.FunctionName
}

// GetNamespace function to get job's namespace
func (m *CronJobTriggerResource) GetNamespace() string {
	return *m.Metadata.Namespace
}

// GetMetadata function to get metadata
func (m *FunctionResourceList) GetMetadata() *metav1.ListMeta {
	return m.Metadata
}

// GetMetadata function to get metadata
func (m *CronJobTriggerResourceList) GetMetadata() *metav1.ListMeta {
	return m.Metadata
}

// GetURL returns function's url
func (f *FunctionResource) GetURL() string {
	return "http://" + f.GetName() + "." + *f.Metadata.Namespace + ":" + strconv.Itoa((*f).Spec.Service.Ports[0].Port)
}

// GetName returns function's url
func (f *FunctionResource) GetName() string {
	return *f.Metadata.Name
}
