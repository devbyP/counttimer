package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

var (
	dirName  = ".counttimer"
	confType = "json"
	confName = "ct"

	envDirPath  = "CT_PATH"
	defaultConf = &Config{
		NotifyMethod: "",
		LineUserID:   "",
	}
)

type Config struct {
	NotifyMethod string `json:"notifyMethod,omitempty"`
	LineUserID   string `json:"lineUserId,omitempty"`
	SaveLog      string `json:"saveLog,omitempty"`
}

func getBasePath() (string, error) {
	var err error
	targetPath := os.Getenv(envDirPath)
	if targetPath == "" {
		targetPath, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(targetPath, dirName), nil
}

func getConf() (*Config, error) {
	basePath, err := getBasePath()
	if err != nil {
		return nil, err
	}
	confFileName := confName + "." + confType
	fullPath := filepath.Join(basePath, confFileName)
	f, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	return conf, nil
}

func initBase() error {
	basepath, err := getBasePath()
	if err != nil {
		return err
	}
	return os.MkdirAll(basepath, os.ModePerm)
}
