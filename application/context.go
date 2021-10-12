package application

import (
	context2 "context"
	"github.com/nekinci/paas/core"
	"github.com/nekinci/paas/specification"
	"golang.org/x/sync/semaphore"
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
	embeddedApplications map[string]core.Host
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

func (context *Context) KillApplication(application Application) {
	id, err := application.Kill()

	if err != nil {
		application.SetStatus(ZOMBIE)
		context.invalidateApplication(application)
		return
	}

	if *id == application.GetApplicationInfo().Id[:6] {
		context.invalidateApplication(application)
		application.SetStatus(STOPPED)
	} else {
		application.SetStatus(ZOMBIE)
		context.invalidateApplication(application)
	}
}

// Get returns a Host.
func (context Context) Get(key string) core.Host {
	// return UrlStrategy(key, context)
	return WildcardStrategy(key, context)
}

func (context Context) addValidApplication(application Application) {
	context.validApplications[application.GetApplicationInfo().Name] = &application
}

func (context Context) invalidateApplication(application Application) {
	context.invalidApplications[application.GetApplicationInfo().Name] = &application
	delete(context.validApplications, application.GetApplicationInfo().Name)
}

// NewContext returns new context with embeddedApplications: core.Host
func NewContext() *Context {
	embeddedApplications := make(map[string]core.Host)
	embeddedApplications[""] = core.NewEmbeddedTcpApplication("", "5000")
	embeddedApplications["www"] = core.NewEmbeddedTcpApplication("www", "5000")
	embeddedApplications["frontend"] = core.NewEmbeddedTcpApplication("frontend", "5100")
	return &Context{
		validApplications:    make(map[string]*Application),
		invalidApplications:  make(map[string]*Application),
		sema:                 semaphore.NewWeighted(int64(applicationLimit)),
		ctx:                  context2.Background(),
		embeddedApplications: embeddedApplications,
	}
}
