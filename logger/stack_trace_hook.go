package logger

import (
	"github.com/facebookgo/stack"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	LOGRUS_STACK_JUMP           = 4
	LOGRUS_FIELDLESS_STACK_JUMP = 6
)

type StackHook struct {
	CallerLevels []logrus.Level
	StackLevels  []logrus.Level
}

func NewStackHook() *StackHook {
	return &StackHook{
		CallerLevels: logrus.AllLevels,
		StackLevels:  []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel},
	}
}

func (h *StackHook) Levels() []logrus.Level {
	return h.CallerLevels
}

func (h *StackHook) Fire(e *logrus.Entry) error {
	skipFrames := LOGRUS_STACK_JUMP
	if len(e.Data) == 0 {
		skipFrames = LOGRUS_FIELDLESS_STACK_JUMP
	}

	var frames stack.Stack

	_frames := stack.Callers(skipFrames)

	for _, frame := range _frames {
		if !strings.Contains(frame.File, "github.com/sirupsen/logrus") {
			frames = append(frames, frame)
		}
	}
	if len(frames) > 0 {
		e.Data["func"] = frames[0].String()
		e.Data["line"] = frames[0].Line
		e.Data["file"] = frames[0].File
		// Set the available frames to "stack" field.
		for _, level := range h.StackLevels {
			if e.Level == level {
				e.Data["errTrace"] = frames
				break
			} else {
				delete(e.Data, "errTrace")
			}
		}
	}

	return nil
}
