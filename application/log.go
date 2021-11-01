package application

import (
	"fmt"
	"time"
)

// LogType is used to classify to log.
type LogType int8

const (
	REMOVE  LogType = 1
	INFO    LogType = 2
	WARNING LogType = 3
	ERROR   LogType = 4
)

// Log struct which using by Application for logging activities.
type Log struct {
	LogType LogType
	Content string
	Show    bool // Should logContent shown to end user?
	LogTime time.Time
}

// NewLog returns new log.
func NewLog(logType LogType, content string, show bool) Log {
	return Log{
		LogType: logType,
		Content: content,
		Show:    show,
		LogTime: time.Now(),
	}
}

// NewInvisibleLog returns new log which not shown to end user.
func NewInvisibleLog(logType LogType, content string) Log {
	return NewLog(logType, content, false)
}

// NewVisibleLog returns new log which shown to end user.
func NewVisibleLog(logType LogType, content string) Log {
	return NewLog(logType, content, true)
}

// NewInfoLog returns new INFO log. The log visible to user as default.
func NewInfoLog(content string) Log {
	return NewLog(INFO, content, true)
}

// NewWarningLog returns new WARNING log. The log visible to user as default.
func NewWarningLog(content string) Log {
	return NewLog(WARNING, content, true)
}

// NewErrorLog returns new ERROR log. The log visible to user as default.
func NewErrorLog(content string) Log {
	return NewLog(ERROR, content, true)
}

// NewRemoveLog returns new REMOVE log. The log visible to user as default.
func NewRemoveLog(content string) Log {
	return NewLog(REMOVE, content, true)
}

type String struct {
	string
}

// FormatString shadows the return type of string to String.
func FormatString(format string, args ...interface{}) String {
	return String{fmt.Sprintf(format, args)}
}

// ToWarningLog converts String to Log.
func (s String) ToWarningLog() Log {
	return NewWarningLog(s.string)
}

// ToErrorLog converts String to Log.
func (s String) ToErrorLog() Log {
	return NewErrorLog(s.string)
}

// ToRemoveLog converts String to Log.
func (s String) ToRemoveLog() Log {
	return NewRemoveLog(s.string)
}

// ToInfoLog converts String to Log.
func (s String) ToInfoLog() Log {
	return NewInfoLog(s.string)
}

// VisibleLog converts String to Log.
func (s String) VisibleLog(logType LogType) Log {
	return NewVisibleLog(logType, s.string)
}

// InvisibleLog converts String to Log.
func (s String) InvisibleLog(logType LogType) Log {
	return NewInvisibleLog(logType, s.string)
}

// LogHandler
type LogHandler func(Log)

func (logType LogType) String() string {
	if logType == REMOVE {
		return "REMOVE"
	} else if logType == WARNING {
		return "WARNING"
	} else if logType == INFO {
		return "INFO"
	} else if logType == ERROR {
		return "ERROR"
	}

	return "UNKNOWN"
}

func (log Log) Format() string {
	return fmt.Sprintf("%s: [%s] -> %s", log.LogType.String(), log.LogTime.Format("02.01.2006 15:04:05"), log.Content)
}
