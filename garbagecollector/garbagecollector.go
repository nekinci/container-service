package garbagecollector

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/nekinci/paas/container"
	"log"
)

const (
	retryLimit = 4
)

func ScheduleCollect(ctx *container.Context) {
	_ = gocron.Every(15).Minute().Do(collect, ctx)
	<-gocron.Start()
}

func collect(ctx *container.Context) {
	log.Printf("Garbage collector running...") // Todo implement log levels..
	containers := ctx.InvalidContainers()

	// TODO Container which has zombie status must kill before remove the image.

	for _, value := range containers {

		if !value.IsRemovable {
			continue
		}

		if value.Status == container.RUNNING {
			continue
		}

		if value.Status == container.WAITING {
			continue
		}

		if value.Status == container.READY {
			continue
		}

		if value.Status == container.PAUSED {
			log.Printf("Container removing from garbage collector: %s", value.Id)
			// As first, we should remove killed container and set status to STOPPED.
			err := value.RemoveContainer()
			if err != nil {
				value.RemoveLogs.Logs = append(value.RemoveLogs.Logs, container.Log(fmt.Sprintf("%v", err)))
				continue
			}
			value.Status = container.STOPPED
			continue
		}

		startTime := value.StartTime

		if startTime.After(startTime.Add(value.CacheTime)) {
			err := value.RemoveImage()
			if err != nil {
				value.RemoveLogs.Logs = append(value.RemoveLogs.Logs, container.Log(fmt.Sprintf("%v", err)))

				value.RemoveLogs.Mutex.Lock()
				value.RemoveLogs.RetryCount += 1
				value.RemoveLogs.Mutex.Unlock()

				if retryLimit == value.RemoveLogs.RetryCount {
					// We should notify system administrator through may be email or telegram message..
				}
			}

		}
	}

	log.Printf("Garbage collector ended...")
}
