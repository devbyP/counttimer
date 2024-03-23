package main

import "fmt"

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

func getProperProgramName(p string) string {
	_, ok := preset[p]
	if !ok {
		return Default
	}
	return p
}

func handleCustom(m, s int) (int, error) {
	count := getCountFromMS(m, s)
	if count == 0 {
		return 0, fmt.Errorf("custom program require set time using m and(or) s flag, count == %d", count)
	}
	return count, nil
}
