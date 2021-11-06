package garbagecollector

import (
	"github.com/jasonlvhit/gocron"
	"github.com/nekinci/paas/application"
	"log"
	"sync"
	"time"
)

const (
	retryLimit = 4
)

type applicationRetry struct {
	mu             sync.Mutex
	applicationMap map[string]int
}

func ScheduleCollect(ctx *application.Context) {
	log.Printf("Garbage collector scheduled..\n")
	retries := applicationRetry{
		applicationMap: map[string]int{},
	}
	_ = gocron.Every(1).Minute().Do(collect, ctx, retries)
	<-gocron.Start()
}

func collect(ctx *application.Context, applicationRetries applicationRetry) {
	log.Printf("Garbage collector running...") // Todo implement log levels..
	applications := ctx.InvalidApplications()

	// TODO Container which has zombie status must kill before remove the image.

	for _, v := range applications {
		value := *v
		//if !value.GetIsRemovable() {
		//	continue
		//}

		if value.GetStatus() == application.RUNNING {
			continue
		}

		if value.GetStatus() == application.WAITING {
			continue
		}

		if value.GetStatus() == application.READY {
			continue
		}

		if value.GetStatus() == application.PAUSED {
			log.Printf("Container removing by garbage collector: %s", value.GetApplicationInfo().Id)
			// As first, we should remove killed container and set status to STOPPED.
			err := value.RemoveApplication()
			if err != nil {
				value.AddNewLog(application.FormatString("Container not removed: %v", err).ToRemoveLog())
				continue
			}
			value.SetStatus(application.STOPPED)
			continue
		}

		startTime, _ := time.Parse(time.RFC3339, value.GetApplicationInfo().StartTime)
		cacheTime := value.GetCacheTime()
		if startTime.After(startTime.Add(cacheTime)) {
			log.Printf("Image cleaned from system... %v\n", value.GetApplicationInfo().Image)

			// We used nginx image as "TryIt". So, we shouldn't remove the image for performance.
			if !ctx.AnyValidApplication(value.GetApplicationInfo().Image) && value.GetApplicationInfo().Image != "nginx" {
				err := value.RemoveFromFileSystem()
				if err != nil {
					value.AddNewLog(application.FormatString("Image/File not removed: %v", err).ToRemoveLog())

					applicationRetries.addRetry(value.GetApplicationInfo().Name)
					if retryLimit == applicationRetries.getRetryCount(value.GetApplicationInfo().Name) {
						// We should notify system administrator through may be email or telegram message..
					}
				}
			}

		}
	}

	log.Printf("Garbage collector ended...")
}

func (applicationRetry *applicationRetry) addRetry(app string) {
	applicationRetry.mu.Lock()
	defer applicationRetry.mu.Unlock()

	if _, ok := applicationRetry.applicationMap[app]; !ok {
		applicationRetry.applicationMap[app] = 0
	}

	applicationRetry.applicationMap[app] += 1
}

func (applicationRetry *applicationRetry) getRetryCount(app string) int {

	if value, ok := applicationRetry.applicationMap[app]; ok {
		return value
	}

	return 0
}
