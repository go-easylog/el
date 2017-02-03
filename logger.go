package el

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	FATAL = iota
	ERROR
	WARN
	INFO
	TRACE
)

var levelTag = []string{
	"[FATAL]",
	"[ERROR]",
	"[WARN]",
	"[INFO]",
	"[TRACE]",
}

type logStatus struct {
	Debug       bool
	LogLevel    LogLevel
	RotateLevel RotateLevel
	Rotator     *Rotator
	chRotate    chan *os.File
	chExit      chan bool
}

var stdRotator = NewStdRotator()

var state = &logStatus{
	Debug:       false,
	LogLevel:    WARN,
	RotateLevel: RNone,
	Rotator:     stdRotator,
	chRotate:    make(chan *os.File),
	chExit:      make(chan bool),
}

var syncMtx = new(sync.Mutex)

func SetLogLevel(level LogLevel) {
	syncMtx.Lock()
	defer syncMtx.Unlock()
	state.LogLevel = level
}

func GetLogLevel() LogLevel {
	syncMtx.Lock()
	defer syncMtx.Unlock()
	return state.LogLevel
}

func SetDebug() {
	syncMtx.Lock()
	defer syncMtx.Unlock()
	state.Debug = true
}
func GetDebug() bool {
	syncMtx.Lock()
	defer syncMtx.Unlock()
	return state.Debug
}

func SetRelease() {
	syncMtx.Lock()
	defer syncMtx.Unlock()
	state.Debug = false
}

func IsDebug() bool {
	syncMtx.Lock()
	defer syncMtx.Unlock()
	return state.Debug
}

func File() *os.File {
	syncMtx.Lock()
	defer syncMtx.Unlock()
	return state.Rotator.File()
}

func SetStdoutLog() {
	syncMtx.Lock()
	defer syncMtx.Unlock()
	state.Rotator.Close()
	state.Rotator = stdRotator
	return
}

func SetRotateLog(path string) error {
	syncMtx.Lock()
	defer syncMtx.Unlock()

	if state.Rotator != nil {
		state.Rotator.Close()
	} else {
		go watch()

	}

	r, err := NewRotator(path)
	if err != nil {
		return err
	}

	state.Rotator = r

	r.SetRotateChannel(r.chRotate)

	return nil
}

func Close() {
	state.chExit <- true
	<-state.chExit
	state.Rotator.Close()
}

func watch() {
Exit:
	for {
		select {
		case <-state.chRotate:

		case <-state.chExit:
			break Exit
		}
	}

	state.chExit <- true
}

func Println(l LogLevel, vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()

	if l > state.LogLevel {
		return
	}
	vals = append([]interface{}{levelTag[l]}, vals...)
	fmt.Fprintln(state.Rotator.File(), vals...)
}

func Printf(l LogLevel, f string, vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()

	if l > state.LogLevel {
		return
	}
	fmt.Fprintf(state.Rotator.File(), levelTag[l]+" "+f, vals...)
}

func Fatal(vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()

	outputln(FATAL, vals)

	os.Exit(-1)
}

func Fatalf(f string, vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()
	outputf(FATAL, f, vals)
	os.Exit(-1)
}

func Error(vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()
	if ERROR > state.LogLevel {
		return
	}

	outputln(ERROR, vals)
}

func Errorf(f string, vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()

	if ERROR > state.LogLevel {
		return
	}
	outputf(ERROR, f, vals)
}

func Warn(vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()

	if WARN > state.LogLevel {
		return
	}
	outputln(WARN, vals)
}

func Warnf(f string, vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()

	if WARN > state.LogLevel {
		return
	}
	outputf(WARN, f, vals)
}

func Info(vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()

	if INFO > state.LogLevel {
		return
	}
	outputln(INFO, vals)
}

func Infof(f string, vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()

	if INFO > state.LogLevel {
		return
	}
	outputf(INFO, f, vals)
}

func Trace(vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()

	if TRACE > state.LogLevel {
		return
	}
	outputln(TRACE, vals)
}

func Tracef(f string, vals ...interface{}) {
	syncMtx.Lock()
	defer syncMtx.Unlock()

	if TRACE > state.LogLevel {
		return
	}
	outputf(TRACE, f, vals)
}

func outputln(level LogLevel, vals []interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	vals = append([]interface{}{now, levelTag[level]}, vals...)
	if state.Debug {
		f := makeFormat(vals)
		fmt.Fprintf(state.Rotator.File(), f+"\n", vals...)
	} else {
		fmt.Fprintln(state.Rotator.File(), vals...)
	}
}

func outputf(level LogLevel, f string, vals []interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(state.Rotator.File(), now+" "+levelTag[level]+f, vals...)
}

func makeFormat(vals []interface{}) string {
	f := "%s"
	for _, v := range vals[1:] {
		if _, ok := v.(error); ok {
			f += " %+v"
		} else {
			f += " %v"
		}
	}
	return f
}
