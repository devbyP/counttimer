package main

import "time"

type DotCircle struct {
	limit   int
	current int
	dots    chan string
}

func (c *DotCircle) runDotGen(done chan struct{}, interval time.Duration) {
	for {
		select {
		case <-done:
			return
		default:
			time.Sleep(interval)
			d := ""
			for i := 0; i < c.current; i++ {
				d += "."
			}
			c.current = (c.current % c.limit) + 1
			c.dots <- d
		}
	}
}

func (c *DotCircle) Close() {
	close(c.dots)
}
