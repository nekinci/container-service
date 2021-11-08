package application

import (
	context2 "context"
	"errors"
	"github.com/nekinci/paas/specification"
	"log"
	"sync"
	"time"
)

const (
	applicationLimit    = 7
	applicationLifeTime = 3
)

// InMemory context.
type Context struct {
	validApplications    map[string]*Application
	invalidApplications  map[string]*Application
	ctx                  context2.Context
	embeddedApplications map[string]Host
	stateListeners       []StateListener
	stateMu              sync.Mutex
	appMu                sync.Mutex
	tryItApps            map[string][]string
}

func (context *Context) RunApplication(specification specification.Specification, removable bool) (*Application, error) {

	newApplication := NewApplication(specification, removable)
	err := newApplication.Run()
	if err != nil {
		return nil, err
	}

	newApplication.SetStatus(RUNNING)
	context.addValidApplication(newApplication)

	return &newApplication, nil
}

func (context *Context) GetApplication(app string) Application {
	application := context.validApplications[app]
	if application == nil {
		return nil
	}
	status := (*application).GetStatus()
	if status == WAITING {
		return nil
	}
	return *application
}

func (context *Context) GetApplicationsByUser(email string) []string {
	var userApps []string = []string{}

	for _, value := range context.validApplications {
		app := *value
		if app.GetApplicationInfo().UserEmail == email {
			userApps = append(userApps, app.GetApplicationInfo().Name)
		}
	}

	for username, appnames := range context.tryItApps {
		if username == email {
			for _, appname := range appnames {
				userApps = append(userApps, appname)
			}
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

}

func (context *Context) ScheduleKill(application *Application) {
	timer := time.NewTimer(applicationLifeTime * time.Minute)
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
	context.appMu.Lock()
	defer context.appMu.Unlock()
	context.validApplications[application.GetApplicationInfo().Name] = &application
	keyCount := getMapSize(context.validApplications)

	for _, f := range context.stateListeners {
		f(NewStateEvent(VALIDATE, "Application started", keyCount))
	}

}

func (context *Context) InvalidApplications() map[string]*Application {
	return context.invalidApplications
}

func (context Context) invalidateApplication(application Application) {
	context.appMu.Lock()
	defer context.appMu.Unlock()
	context.invalidApplications[application.GetApplicationInfo().Name] = &application
	delete(context.validApplications, application.GetApplicationInfo().Name)

	keyCount := getMapSize(context.validApplications)

	for _, f := range context.stateListeners {
		f(NewStateEvent(INVALIDATE, "Application stopped.", keyCount))
	}
}

func (context *Context) Handle(app *specification.Specification, removable bool) error {
	context.appMu.Lock()
	appCount := getMapSize(context.validApplications)
	isThereName := context.AnyValidApplicationByName(app.Name)
	context.appMu.Unlock()

	if appCount >= applicationLimit+1 {
		return errors.New("Application limit exceeded\n")
	}

	if isThereName {
		return errors.New("Application name already taken\n")
	}

	if isReservedName(app.Name) {
		return errors.New("This name reserved\n")
	}

	application, err := context.RunApplication(*app, removable)
	if err != nil {
		return err
	}

	if app.Image == "nginx" && app.Name == "nginxapp" {
		log.Printf("Not scheduled to kill for nginxapp")
		return nil
	}
	go context.ScheduleKill(application)
	return nil
}

func (context *Context) AddTryItUser(appName, userName string) {
	if _, ok := context.tryItApps[userName]; !ok {
		context.tryItApps[userName] = []string{}
	}

	anyApp := func() bool {
		dataList := context.tryItApps[userName]
		for _, v := range dataList {
			if v == appName {
				return true
			}
		}
		return false
	}

	if anyApp() {
		return
	}

	context.tryItApps[userName] = append(context.tryItApps[userName], appName)

	removeFunc := func(userName, appName string) {

		dataList := context.tryItApps[userName]
		removeIndex := -1
		for i := range dataList {
			if dataList[i] == appName {
				removeIndex = i
			}
		}

		newDataList := append(dataList[:removeIndex], dataList[removeIndex+1:]...)
		context.tryItApps[userName] = newDataList

	}

	go func() {
		timer := time.NewTimer(applicationLifeTime * time.Minute)
		done := make(chan bool)
		go func() {
			<-timer.C
			done <- true
		}()
		<-done
		removeFunc(userName, appName)
	}()

}

func (context Context) AnyValidApplicationByImage(image string) bool {
	for _, v := range context.validApplications {
		value := *v
		if value.GetApplicationInfo().Image == image {
			return true
		}
	}
	return false
}

func (context Context) AnyValidApplicationByName(appName string) bool {
	for k, _ := range context.validApplications {
		if k == appName {
			return true
		}
	}
	return false
}

func getMapSize(dict map[string]*Application) int {
	var count int
	for _, _ = range dict {
		count = count + 1
	}
	return count
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
		ctx:                  context2.Background(),
		embeddedApplications: embeddedApplications,
		stateListeners:       []StateListener{},
		tryItApps:            make(map[string][]string),
	}
}
