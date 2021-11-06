package application

import "time"

type StateType string

const (
	VALIDATE   StateType = "VALIDATE"
	INVALIDATE StateType = "INVALIDATE"
)

type StateEvent struct {
	Type                StateType
	Content             string
	ApplicationCount    int
	Time                time.Time
	MaxApplicationCount int
}

type StateListener func(StateEvent)

func (context *Context) AddStateListener(listener StateListener) {
	context.stateMu.Lock()
	defer context.stateMu.Unlock()
	context.stateListeners = append(context.stateListeners, listener)
}

func NewStateEvent(stateType StateType, content string, appCount int) StateEvent {
	return StateEvent{
		Type:                stateType,
		Content:             content,
		ApplicationCount:    appCount,
		Time:                time.Now(),
		MaxApplicationCount: applicationLimit,
	}
}
