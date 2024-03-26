package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"
)

var (
	logFileName = "ctlog"
	logDirName  = "logs"
)

type TimeLog struct {
	StartTime   time.Time
	EndTime     time.Time
	Program     string
	Title       string
	Description string
	EarlyFinish bool
}

func (tl TimeLog) writeFormat() string {
	return "test v2\n"
}

func saveLog(tl TimeLog) error {
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
	_, err = f.WriteString(tl.writeFormat())
	if err != nil {
		return err
	}
	return nil
}

func readRecentLog() (string, error) {
	basePath, err := getBasePath()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(basePath, logFileName)
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	info, err := f.Stat()
	if err != nil {
		return "", err
	}
	predictSize := info.Size() - 2048
	start := int64(math.Max(0, float64(predictSize)))
	data := []byte{}
	_, err = f.ReadAt(data, start)
	if err != nil {
		return "", err
	}

	return "", nil
}

func getDateYM(t time.Time) string {
	month := t.Month().String()
	year := t.Year()
	return fmt.Sprintf("-%d-%s", year, month)
}

func initLog() error {
	basePath, err := getBasePath()
	if err != nil {
		return err
	}
	logpath := filepath.Join(basePath, logDirName)
	return os.MkdirAll(logpath, os.ModePerm)
}
