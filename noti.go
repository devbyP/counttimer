package main

import (
// "net/http"
)

type NotifyOption struct {
	Method string
	User   string
	Title  string
}

func NewOptionFromConfig(c *Config, title string) *NotifyOption {
	return &NotifyOption{
		User:   c.LineUserID,
		Method: c.NotifyMethod,
		Title:  title,
	}
}

func Notify(opt *NotifyOption) error {
	return nil
}
