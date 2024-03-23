package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	c   bool
	pre string
	m   int
	s   int
)

func init() {
	flag.BoolVar(&c, "c", false, "custom mode")
	flag.StringVar(&pre, "p", Default, "preset program")
	flag.IntVar(&m, "m", 0, "stop time in minute (custom program)")
	flag.IntVar(&s, "s", 0, "stop time in second (custom program)")
	// d := flag.String("d", "", "count session description for notification")
	// n := flag.Bool("n", true, "notify service")
	customUsage()
	flag.Parse()
}

func main() {
	err := initBase()
	if err != nil {
		printErrHelp(err)
		return
	}
	err = validateFlag(pre, c, m, s)
	if err != nil {
		printErrHelp(err)
		return
	}
	conf, err := getConf()
	if err != nil {
		perr := &os.PathError{}
		if !errors.As(err, &perr) {
			printErrHelp(err)
			return
		}
		conf = defaultConf
	}
	if conf.NotifyMethod != "" {
		fmt.Println("do notify")
	}
	if c {
		pre = Custom
	}
	program := getProperProgramName(pre)

	countLimit, err := getCountBaseOnProgram(program, m, s)
	if err != nil {
		printErrHelp(err)
		return
	}

	printIntro(program, countLimit)
	err = saveLog(TimeLog{})
	if err != nil {
		printErrHelp(err)
		return
	}

	startTime := time.Now()

	done := make(chan struct{})
	go gracefulShutdown(done)
	dc := &DotCircle{
		limit:   3,
		current: 1,
		dots:    make(chan string),
	}
	defer func() {
		close(done)
		dc.Close()
	}()
	go dc.runDotGen(done, time.Millisecond*750)
	go countTimer(done, countLimit, dc, startTime)

	<-done
	printEnding()
}

func validateFlag(pre string, c bool, m, s int) error {
	if pre != Default && c {
		return fmt.Errorf("cannot decided program")
	}
	if c && (m <= 0 && s <= 0) {
		return fmt.Errorf("custom flag require -m and(or) -s flag")
	}
	return nil
}

func customUsage() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println("\npreset program value list:")
		fmt.Println(`
  default
  longbreak
  shortbreak
  break    - 10 min break
  longwork - 1 hour
  tea      - 3 minute counting for steeping
  b        - alias for break
  lb       - alias for longbreak
  sb       - alias for shortbreak
  inf      - program counting infinite until user stop the program
  custom   - for custom program (require -m and[or] -s)`)
	}
}

func gracefulShutdown(done chan struct{}) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()
	fmt.Println("\ngot interuption signal")
	done <- struct{}{}
}

func printErrHelp(err error) {
	fmt.Println(err)
	fmt.Println()
	flag.Usage()
}

func printIntro(program string, c int) {
	sm, ss := getTimeFromCount(c)
	a, ok := aliasPro[program]
	if ok {
		program = a
	}
	fmt.Printf("start count timer \"%s\" program.\n", program)
	fmt.Printf("counting time %02d minute %02d second\n", sm, ss)
}

func printEnding() {
	fmt.Println()
	fmt.Println("done counting")
	fmt.Println(timeFormat(time.Now()))
}

func timeFormat(t time.Time) string {
	s := t.Second()
	m := t.Minute()
	h := t.Hour()
	wday := t.Weekday().String()
	day := t.Day()
	month := t.Month().String()
	return fmt.Sprintf("%s %d %s - %d:%02d:%02d", wday, day, month, h, m, s)
}
