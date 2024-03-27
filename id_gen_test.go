package main

import (
	"testing"
	"time"
)

func Test_genSessionID(t *testing.T) {
	t.Log(genSessionID(time.Now()))
}
