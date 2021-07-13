package untraceable

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/josedelrio85/voalarm"
)

// ErrorLogger is a struct to handle error properties
type ErrorLogger struct {
	Msg    string
	Status int
	Err    error
	Log    string
}

// SendAlarm to VictorOps plattform and format the error for more info
func (e *ErrorLogger) SendAlarm() {
	e.Msg = fmt.Sprintf("Untraceable -> %s %v", e.Msg, e.Err)
	log.Println(e.Log)

	mstype := voalarm.Acknowledgement
	switch e.Status {
	case http.StatusInternalServerError:
		mstype = voalarm.Warning
	case http.StatusUnprocessableEntity:
		mstype = voalarm.Info
	}

	alarm := voalarm.NewClient("")
	_, err := alarm.SendAlarm(e.Msg, mstype, e.Err)
	if err != nil {
		log.Fatalf(e.Msg)
	}
}

// LogError obtains a trace of the line and file where the error happens
func LogError(err error) string {
	pc, fn, line, _ := runtime.Caller(1)
	return fmt.Sprintf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
}
