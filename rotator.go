package el

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RotateLevel int

const (
	RDay = iota
	RMonth
	RYear
	RNone
)

type Rotator struct {
	PathFormat  string
	file        *os.File
	rotateLevel RotateLevel
	chRotate    chan<- *os.File
	chExit      chan bool
}

func NewRotator(path string) (*Rotator, error) {

	var rLevel RotateLevel = RNone

	if strings.Index(path, "%D") != -1 {
		rLevel = RDay
	} else if strings.Index(path, "%M") != -1 {
		rLevel = RMonth
	} else if strings.Index(path, "%Y") != -1 {
		rLevel = RYear
	}

	result := &Rotator{
		PathFormat:  path,
		rotateLevel: rLevel,
		chRotate:    nil,
		chExit:      make(chan bool),
	}

	now := time.Now()
	if err := result.changeLogFile(now); err != nil {
		return nil, err
	}

	go result.watch(&now)

	return result, nil
}

func NewStdRotator() *Rotator {
	result := &Rotator{
		PathFormat:  "",
		rotateLevel: RNone,
		file:        os.Stdout,
		chRotate:    nil,
		chExit:      nil,
	}

	return result
}

func (r *Rotator) File() *os.File {
	return r.file
}

func (r *Rotator) SetRotateChannel(ch chan<- *os.File) {
	r.chRotate = ch
}

func (r *Rotator) changeLogFile(t time.Time) error {

	// open new log
	filePath := logFilePath(r.PathFormat, t)
	fullPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(fullPath)
	// make dir

	if s, err := os.Stat(dirPath); err != nil {
		if err := os.MkdirAll(dirPath, 0775); err != nil {
			return err
		}
	} else if !s.IsDir() {
		return fmt.Errorf("file is not directory")
	}

	lf, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	if r.chRotate != nil {
		r.chRotate <- lf
	}

	if r.file != nil {
		r.file.Close()
	}

	r.file = lf

	return nil
}

func (r *Rotator) watch(lastDay *time.Time) {

	ticker := time.NewTicker(time.Second)

Exit:
	for {
		select {
		case <-ticker.C:
			t := time.Now()

			switch r.rotateLevel {
			case RDay:
				if lastDay.Day() == t.Day() {
					continue
				}
			case RMonth:
				if lastDay.Month() == t.Month() {
					continue
				}
			case RYear:
				if lastDay.Year() == t.Year() {
					continue
				}
			}
			if r.changeLogFile(t) != nil {
				Error("Failed rotate log file")
			}
			lastDay = &t
		case <-r.chExit:
			break Exit
		}
	}
	ticker.Stop()
	r.chExit <- true
}

func (r *Rotator) Close() {
	if r.chExit != nil {
		r.chExit <- true
		<-r.chExit
	}
}

func logFilePath(path string, t time.Time) string {
	path = strings.Replace(path, "%D", fmt.Sprintf("%02d", t.Day()), -1)
	path = strings.Replace(path, "%M", fmt.Sprintf("%02d", t.Month()), -1)
	path = strings.Replace(path, "%Y", fmt.Sprintf("%02d", t.Year()), -1)
	return path
}
