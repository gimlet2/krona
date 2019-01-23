package main

import (
	"log"
	"net/http"

	"github.com/robfig/cron"
)

// Cron interface wrapper around cron lib
type Cron interface {
	Cancel(name string)
	Schedule(name string, url string, pattern string)
	Cleanup(names Set)
	NeedsUpdate(name string, pattern string) bool
	Has(name string) bool
}

// CronImpl default implemetation of Cron interface
type CronImpl map[string]*CronJob

// NewCron constructor of Cron
func NewCron() Cron {
	c := CronImpl{}
	return &c
}

// Cancel cancels job by name
func (c CronImpl) Cancel(name string) {
	c[name].Cron.Stop()
	delete(c, name)
}

// Schedule schedules job for triggering
func (c CronImpl) Schedule(name string, url string, pattern string) {
	log.Printf("Function '%s' with schedule '%s' was descovered", name, pattern)
	cron := cron.New()
	cron.AddFunc(pattern, func() {
		log.Printf("Trigger '%s' function - GET - %s", name, url)
		http.Get(url)
	})
	c[name] = &CronJob{pattern, cron}
	cron.Start()
}

// Cleanup cancels outdated jobs
func (c CronImpl) Cleanup(names Set) {
	for k := range c {
		if names.Has(k) {
			c.Cancel(k)
		}
	}
}

// NeedsUpdate checks was jobs pattern changed or not
func (c CronImpl) NeedsUpdate(name string, pattern string) bool {
	return c[name] != nil && c[name].Pattern != pattern
}

// Has check that job exists
func (c CronImpl) Has(name string) bool {
	return c[name] != nil
}
