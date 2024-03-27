package main

import (
	"encoding/base64"
	"strconv"
	"time"
)

func genSessionID(t time.Time) string {
	d := strconv.FormatInt(t.Unix(), 32)
	tStr := []byte(d)
	return base64.StdEncoding.WithPadding(base64.NoPadding).Strict().EncodeToString(tStr)
}
