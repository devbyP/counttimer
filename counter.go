package main

import (
	"fmt"
	"time"
)

func countTimer(done chan struct{}, countLimit int, dc *DotCircle, start time.Time) {
	if dc == nil {
		// default dot circle
		dc = &DotCircle{
			limit:   3,
			current: 1,
			dots:    make(chan string),
		}
		go dc.runDotGen(done, time.Millisecond*750)
	}
	ticker := time.NewTicker(time.Second * 1)
	count := 0
	locdots := ""
	fmt.Println("timer start at", timeFormat(start))
	fmt.Println()
	for {
		// 0 case for continue count timer
		if countLimit > 0 && countLimit <= count {
			ticker.Stop()
			m, s := getTimeFromCount(count)
			fmt.Printf("\r%02d:%02d     \n", m, s)
			done <- struct{}{}
		}
		select {
		case <-ticker.C:
			count++
			mm, ss := getTimeFromCount(count)
			fmt.Printf("\r%02d:%02d %-4s", mm, ss, locdots)
		case locdots = <-dc.dots:
			mm, ss := getTimeFromCount(count)
			fmt.Printf("\r%02d:%02d %-4s", mm, ss, locdots)
		case <-done:
			return
		}
	}
}

func getCountBaseOnProgram(p string, m, s int) (int, error) {
	pt, ok := preset[p]
	if p == Custom {
		return handleCustom(m, s)
	}
	if !ok {
		p = Default
		pt = preset[p]
	}
	minute := pt[0]
	second := pt[1]
	return getCountFromMS(minute, second), nil
}

func getTimeFromCount(c int) (int, int) {
	m := c / 60
	s := c % 60
	return m, s
}

func getCountFromMS(m, s int) int {
	return s + (m * 60)
}
