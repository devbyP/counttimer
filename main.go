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

type FlagsConf struct {
	IsCustom    bool
	Program     Program
	Title       string
	Description string
	Minute      int
	Second      int
	Notify      bool
}

func handleFlags(f *FlagsConf) {
	var pro string
	flag.BoolVar(&f.IsCustom, "c", false, "custom mode")
	flag.StringVar(&pro, "p", Default, "preset program")
	flag.IntVar(&f.Minute, "m", 0, "stop time in minute (custom program)")
	flag.IntVar(&f.Second, "s", 0, "stop time in second (custom program)")
	flag.StringVar(&f.Title, "t", "test default", "title message")
	flag.StringVar(&f.Description, "d", "", "count session description for notification")
	flag.BoolVar(&f.Notify, "n", true, "notify service after time over")
	customUsage()
	flag.Parse()
	f.Program = Program(pro)
}

func main() {
	f := &FlagsConf{}
	handleFlags(f)
	err := initBase()
	if err != nil {
		printErrHelp(err)
		return
	}
	err = f.validateFlag()
	if err != nil {
		printErrHelp(err)
		return
	}
	f.Default()
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
	program := f.Program.getProperName()

	countLimit, err := f.Program.getCount(f.Minute, f.Second)
	if err != nil {
		printErrHelp(err)
		return
	}
	var startTime, endTime time.Time

	printIntro(program, countLimit)
	startTime = time.Now()
	tl := TimeLog{
		StartTime:   startTime,
		Program:     program,
		Title:       f.Title,
		Description: f.Description,
	}

	early := make(chan struct{})
	done := make(chan struct{})
	go gracefulShutdown(early)
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

	select {
	case <-done:
		endTime = time.Now()
		tl.EndTime = endTime
		// handle logic properly done.
	case <-early:
		endTime = time.Now()
		tl.EndTime = endTime
		tl.EarlyFinish = true
		// handle logic early done.
	}
	err = saveLog(tl)
	if err != nil {
		printErrHelp(err)
		return
	}
	printEnding(endTime)
}

func (f *FlagsConf) validateFlag() error {
	if f.Program != Default && f.IsCustom {
		return fmt.Errorf("cannot decided program")
	}
	if f.IsCustom && (f.Minute <= 0 && f.Second <= 0) {
		return fmt.Errorf("custom flag require -m and(or) -s flag")
	}
	return nil
}

func (f *FlagsConf) Default() {
	if f.IsCustom {
		f.Program = Custom
	}
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

func gracefulShutdown(early chan struct{}) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()
	early <- struct{}{}
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

func printEnding(endTime time.Time) {
	fmt.Println()
	fmt.Println("done counting")
	fmt.Println(timeFormat(endTime))
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
