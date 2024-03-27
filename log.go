package main

import (
	"errors"
	"fmt"
	//"math"
	"os"
	"path/filepath"
	"time"
)

var (
	logFileName = "ctlog"
	logDirName  = "logs"
)

// log mode
var (
	Start = 1
	End   = 2
)

type TimeLog struct {
	SessionID   string
	StartTime   time.Time
	EndTime     time.Time
	Program     string
	Minute      int
	Second      int
	Title       string
	Description string
	EarlyFinish bool
}

func (tl TimeLog) writeFormat(w *os.File, status int) error {
	mode := []string{"Start", "End", "Early End", "Error"}
	var err error
	logFormat := "[%s]{{%s}} ID: %s - %s=%02d:%02d (Program: %s) - Title: %s - Desc: %s\n"
	errFormat := "[%s]{%s} ID: %s - Message: %s"
	switch status {
	case Start:
		f := fmt.Sprintf(
			logFormat,
			mode[0],
			timeFormat(tl.StartTime),
			tl.SessionID,
			"time set",
			tl.Minute,
			tl.Second,
			tl.Program,
			tl.Title,
			tl.Description,
		)
		_, err = w.WriteString(f)
	case End:
		smode := 1
		actualM := tl.Minute
		actualS := tl.Second
		if tl.EarlyFinish {
			smode = 2
			ac := tl.EndTime.Unix() - tl.StartTime.Unix()
			actualM, actualS = getTimeFromCount(int(ac))
		}
		f := fmt.Sprintf(
			logFormat,
			mode[smode],
			timeFormat(tl.EndTime),
			tl.SessionID,
			"actual",
			actualM,
			actualS,
			tl.Program,
			tl.Title,
			tl.Description,
		)
		_, err = w.WriteString(f)
	}
	if err != nil {
		f := fmt.Sprintf(
			errFormat,
			mode[3],
			timeFormat(tl.StartTime),
			tl.SessionID,
			err.Error(),
		)
		_, werr := w.WriteString(f)
		if werr != nil {
			return errors.Join(err, werr)
		}
	}
	return err
}

func saveLog(tl TimeLog, status int) error {
	err := initLog()
	if err != nil {
		return err
	}
	basePath, err := getBasePath()
	if err != nil {
		return err
	}
	basePath = filepath.Join(basePath, logDirName, logFileName+getDateYM(tl.StartTime))
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(basePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	err = tl.writeFormat(f, status)
	if err != nil {
		return err
	}
	return nil
}

// func readRecentLog() (string, error) {
// 	basePath, err := getBasePath()
// 	if err != nil {
// 		return "", err
// 	}
// 	filePath := filepath.Join(basePath, logFileName)
// 	f, err := os.Open(filePath)
// 	if err != nil {
// 		return "", err
// 	}
// 	info, err := f.Stat()
// 	if err != nil {
// 		return "", err
// 	}
// 	predictSize := info.Size() - 2048
// 	start := int64(math.Max(0, float64(predictSize)))
// 	data := []byte{}
// 	_, err = f.ReadAt(data, start)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return "", nil
// }

func getDateYM(t time.Time) string {
	month := t.Month() + 1
	year := t.Year()
	return fmt.Sprintf("-%d-%d", year, month)
}

func initLog() error {
	basePath, err := getBasePath()
	if err != nil {
		return err
	}
	logpath := filepath.Join(basePath, logDirName)
	return os.MkdirAll(logpath, os.ModePerm)
}
