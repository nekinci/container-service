package container

import (
	context2 "context"
	"golang.org/x/sync/semaphore"
	"log"
)

type Context struct {
	runningContainers map[string]*Container
	stoppedContainers map[string]*Container
	queue *semaphore.Weighted
	ctx	context2.Context
}


func NewContext() *Context {
	return &Context{
		runningContainers: make(map[string]*Container),
		stoppedContainers: make(map[string]*Container),
		queue:             semaphore.NewWeighted(2),
		ctx: context2.Background(),
	}
}


func (context *Context) Acquire(container *Container){
	log.Printf("Semaphore: acquiring: %s", container.Id)
	context.queue.Acquire(context.ctx, 1)
	context.runningContainers[container.Specification.Name] = container
}

func (context *Context) Release(container *Container) {
	log.Printf("Semaphore: releasing: %s", container.Id)
	context.queue.Release(1)
	delete(context.runningContainers, container.Specification.Name)
	context.stoppedContainers[container.Specification.Name] = container
}

func (context *Context) Get(key string) *Container {
	return context.runningContainers[key]
}


func (context *Context) InvalidContainers() map[string]*Container {
	return context.stoppedContainers
}