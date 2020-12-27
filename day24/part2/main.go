// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "Verbose log messages")
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		process(arg)
	}

	log.Printf("Done.")
}

func process(filename string) {
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	lines := strings.Split(string(buf), "\n")
	tiles := map[key]bool{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		key := parse(line)
		if tiles[key] {
			delete(tiles, key)
		} else {
			tiles[key] = true
		}
	}

	for day := 1; day <= 100; day++ {
		eod := runDay(tiles)
		log.Printf("Day %v: %v", day, len(eod))
		tiles = eod
	}

	log.Printf("Solution: %v", len(tiles))
}

type key struct {
	x, y int
}

func (k key) e() key { return key{x: k.x + 1, y: k.y} }
func (k key) w() key { return key{x: k.x - 1, y: k.y} }
func (k key) se() key {
	if evenRow(k) {
		return key{x: k.x, y: k.y - 1}
	}
	return key{x: k.x + 1, y: k.y - 1}
}
func (k key) sw() key {
	if evenRow(k) {
		return key{x: k.x - 1, y: k.y - 1}
	}
	return key{x: k.x, y: k.y - 1}
}
func (k key) ne() key {
	if evenRow(k) {
		return key{x: k.x, y: k.y + 1}
	}
	return key{x: k.x + 1, y: k.y + 1}
}
func (k key) nw() key {
	if evenRow(k) {
		return key{x: k.x - 1, y: k.y + 1}
	}
	return key{x: k.x, y: k.y + 1}
}

func runDay(black map[key]bool) map[key]bool {
	eod := map[key]bool{}
	white := map[key]bool{}

	for k := range black {
		bn, wn := neighbors(k, black)
		if bn > 0 && bn <= 2 {
			eod[k] = true
		}
		for k2 := range wn {
			white[k2] = true
		}
	}
	for k := range white {
		bn, _ := neighbors(k, black)
		if bn == 2 {
			eod[k] = true
		}
	}
	return eod
}

func neighbors(k key, black map[key]bool) (int, map[key]bool) {
	var bn int
	white := map[key]bool{}
	f := func(k2 key) {
		if black[k2] {
			bn++
		} else {
			white[k2] = true
		}
	}
	f(k.e())
	f(k.w())
	f(k.se())
	f(k.sw())
	f(k.ne())
	f(k.nw())
	return bn, white
}

func evenRow(k key) bool {
	v := k.y
	if v < 0 {
		v = -v
	}
	return v%2 == 0
}

func oddRow(k key) bool {
	v := k.y
	if v < 0 {
		v = -v
	}
	return v%2 == 1
}

func parse(s string) key {
	var k key

	for len(s) > 0 {
		switch {
		case strings.HasPrefix(s, "w"):
			k.x--
			s = s[1:]
		case strings.HasPrefix(s, "e"):
			k.x++
			s = s[1:]
		case strings.HasPrefix(s, "se") && evenRow(k):
			k.y--
			s = s[2:]
		case strings.HasPrefix(s, "se") && oddRow(k):
			k.x++
			k.y--
			s = s[2:]
		case strings.HasPrefix(s, "sw") && evenRow(k):
			k.x--
			k.y--
			s = s[2:]
		case strings.HasPrefix(s, "sw") && oddRow(k):
			k.y--
			s = s[2:]
		case strings.HasPrefix(s, "ne") && evenRow(k):
			k.y++
			s = s[2:]
		case strings.HasPrefix(s, "ne") && oddRow(k):
			k.x++
			k.y++
			s = s[2:]
		case strings.HasPrefix(s, "nw") && evenRow(k):
			k.x--
			k.y++
			s = s[2:]
		case strings.HasPrefix(s, "nw") && oddRow(k):
			k.y++
			s = s[2:]
		default:
			log.Fatalf("parse error: %v", s)
		}
	}
	return k
}

func check(fmtStr string, args ...interface{}) {
	if err := args[len(args)-1]; err != nil {
		log.Fatalf(fmtStr, args...)
	}
}

func logf(fmt string, args ...interface{}) {
	if *verbose {
		log.Printf(fmt, args...)
	}
}
