package log

import(
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type StacktraceHook struct {
	innerHook logrus.Hook
}

func (h *StacktraceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *StacktraceHook) Fire(e *logrus.Entry) error {
	if v, found := e.Data[logrus.ErrorKey]; found {
		if err, iserr := v.(error); iserr {
			type stackTracer interface {
				StackTrace() errors.StackTrace
			}
			if st, isst := err.(stackTracer); isst {
				stack := fmt.Sprintf("%+v", st.StackTrace())
				e.Data["stacktrace"] = stack
			}
		}
	}
	h.innerHook.Fire(e)
	return nil
}