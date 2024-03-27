package main

import "fmt"

type Program string

const (
	Default = "default"
	Custom  = "custom"
	Infinte = "inf"
)

var preset = map[string][2]int{
	Infinte:      {0, 0},
	"longWork":   {60, 0},
	Default:      {25, 0},
	"break":      {10, 0},
	"b":          {10, 0},
	"shortbreak": {5, 0},
	"sb":         {5, 0},
	"longbreak":  {15, 0},
	"lb":         {15, 0},
	Custom:       {0, 0},
	"tea":        {3, 0},
}

var aliasPro = map[string]string{
	"sb": "shortbreak",
	"b":  "break",
	"lb": "longbreak",
}

func (p Program) getProperName() string {
	name := string(p)
	_, ok := preset[name]
	if !ok {
		return Default
	}
	if a, ok := aliasPro[name]; ok {
		return a
	}
	return name
}

func (p Program) getTime(m, s int) (int, int) {
	if string(p) == Custom {
		return m, s
	}
	pt, ok := preset[string(p)]
	if !ok {
		pt = preset[Default]
	}
	return pt[0], pt[1]
}

func (p Program) getCount(m, s int) (int, error) {
	pt, ok := preset[string(p)]
	if string(p) == Custom {
		return p.handleCustom(m, s)
	}
	if !ok {
		p = Default
		pt = preset[string(p)]
	}
	minute := pt[0]
	second := pt[1]
	return getCountFromMS(minute, second), nil
}

func (p Program) handleCustom(m, s int) (int, error) {
	count := getCountFromMS(m, s)
	if count == 0 {
		return 0, fmt.Errorf("custom program require set time using m and(or) s flag, count == %d", count)
	}
	return count, nil
}
