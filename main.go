package main

import (
	"log"
	"os"
	"time"
)

func main() {
	log.Print("Start")
	client := NewClient(os.Getenv("KUBE_CONFIG"))
	cron := NewCron()

	for {
		cronTriggers := client.GetCronTriggers()
		for _, cronTrigger := range cronTriggers.Items {
			function, err := client.GetFunction(cronTrigger.GetNamespace(), cronTrigger.GetFunctionName())
			if err != nil {
				continue
			}
			if cron.NeedsUpdate(function.GetName(), cronTrigger.GetPattern()) {
				cron.Cancel(function.GetName())
			}
			if !cron.Has(function.GetName()) {
				cron.Schedule(function.GetName(), function.GetURL(), cronTrigger.GetPattern())
			}

		}
		cleanup(&cron, cronTriggers)
		time.Sleep(5 * time.Second)
	}
}

func cleanup(cron *Cron, cronTriggers CronJobTriggerResourceList) {
	var names []string
	for _, cronTrigger := range cronTriggers.Items {
		names = append(names, cronTrigger.GetFunctionName())
	}
	(*cron).Cleanup(NewSet(names))
}
