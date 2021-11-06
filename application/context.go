package application

import (
	context2 "context"
	"github.com/nekinci/paas/specification"
	"golang.org/x/sync/semaphore"
	"time"
)

const (
	applicationLimit = 2
)

// InMemory context.
type Context struct {
	validApplications    map[string]*Application
	invalidApplications  map[string]*Application
	sema                 *semaphore.Weighted
	ctx                  context2.Context
	embeddedApplications map[string]Host
}

func (context *Context) RunApplication(specification specification.Specification) *Application {

	context.sema.Acquire(context.ctx, 1)

	newApplication := NewApplication(specification)
	err := newApplication.Run()
	if err != nil {
		context.sema.Release(1)
		return nil
	}

	newApplication.SetStatus(RUNNING)
	context.addValidApplication(newApplication)

	return &newApplication
}

func (context *Context) GetApplication(app string) Application {
	application := context.Get(app).(Application)
	return application
}

func (context *Context) GetApplicationsByUser(email string) []string {
	var userApps []string = []string{}

	for _, value := range context.validApplications {
		app := *value
		if app.GetApplicationInfo().UserEmail == email {
			userApps = append(userApps, app.GetApplicationInfo().Name)
		}
	}

	return userApps
}

func (context *Context) KillApplication(application Application) {
	id, err := application.Kill()

	if err != nil {
		application.SetStatus(ZOMBIE)
		context.invalidateApplication(application)
		return
	}

	if *id == application.GetApplicationInfo().Id[:6] {
		context.invalidateApplication(application)
		application.SetStatus(PAUSED)
	} else {
		application.SetStatus(ZOMBIE)
		context.invalidateApplication(application)
	}

	context.sema.Release(1)

}

func (context *Context) ScheduleKill(application *Application) {
	timer := time.NewTimer(1 * time.Minute)
	done := make(chan bool)
	go func() {
		<-timer.C
		done <- true
	}()
	<-done
	context.KillApplication(*application)
}

// Get returns a Host.
func (context Context) Get(key string) Host {
	// return UrlStrategy(key, context)
	return WildcardStrategy(key, context)
}

func (context Context) addValidApplication(application Application) {
	context.validApplications[application.GetApplicationInfo().Name] = &application
}

func (context *Context) InvalidApplications() map[string]*Application {
	return context.invalidApplications
}

func (context Context) invalidateApplication(application Application) {
	context.invalidApplications[application.GetApplicationInfo().Name] = &application
	delete(context.validApplications, application.GetApplicationInfo().Name)
}

func (context Context) AnyValidApplication(image string) bool {
	for _, v := range context.validApplications {
		value := *v
		if value.GetApplicationInfo().Image == image {
			return true
		}
	}
	return false
}

// NewContext returns new context with embeddedApplications: core.Host
func NewContext() *Context {
	embeddedApplications := make(map[string]Host)
	embeddedApplications[""] = NewEmbeddedTcpApplication("", "3000")
	embeddedApplications["www"] = NewEmbeddedTcpApplication("www", "3000")
	embeddedApplications["frontend"] = NewEmbeddedTcpApplication("frontend", "5200")
	embeddedApplications["api"] = NewEmbeddedTcpApplication("api", "8070")
	return &Context{
		validApplications:    make(map[string]*Application),
		invalidApplications:  make(map[string]*Application),
		sema:                 semaphore.NewWeighted(int64(applicationLimit)),
		ctx:                  context2.Background(),
		embeddedApplications: embeddedApplications,
	}
}
