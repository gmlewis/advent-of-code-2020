// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"fmt"
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

	cups := strings.TrimSpace(string(buf))
	// log.Printf("cups: %v", cups)
	game := &Game{}
	for _, cup := range cups {
		game.cups = append(game.cups, cup)
	}

	for i := 0; i < 100; i++ {
		game.move()
	}
	log.Printf("game=%v", game)

	log.Printf("Solution: %v", game.solution())
}

type Game struct {
	cups   []rune
	curCup int
}

func (g *Game) move() {
	destLabel := g.cups[g.curCup]
	bumpDest := func() {
		// log.Printf("BEFORE bumpDest: destLabel=%c", destLabel)
		destLabel = ((destLabel - '1' + rune(len(g.cups)-1)) % rune(len(g.cups))) + '1'
		// log.Printf("AFTER bumpDest: destLabel=%c", destLabel)
	}
	bumpDest()

	newOrder := []rune{g.cups[g.curCup]}
	var pickUp []rune
	avoid := map[rune]bool{}
	for i := 1; i <= 3; i++ {
		cup := g.cups[(g.curCup+i)%len(g.cups)]
		avoid[cup] = true
		pickUp = append(pickUp, cup)
	}
	for avoid[destLabel] {
		bumpDest()
	}

	for i := 4; i < len(g.cups); i++ {
		cup := g.cups[(g.curCup+i)%len(g.cups)]
		newOrder = append(newOrder, cup)
		if destLabel == cup {
			newOrder = append(newOrder, pickUp...)
		}
	}
	// log.Printf("cups=%v, pickUp=%v, newOrder=%v, destLabel=%c", g, dump(pickUp), dump(newOrder), destLabel)
	g.curCup = 1
	g.cups = newOrder
}

func (g *Game) solution() string {
	var r string
	var i int
	for ; g.cups[i] != '1'; i++ {
	}
	for i++; len(r)+1 < len(g.cups); i++ {
		r += fmt.Sprintf("%c", g.cups[i%len(g.cups)])
	}
	return r
}

func (g *Game) String() string {
	var s string
	for i, c := range g.cups {
		if i == g.curCup {
			s += fmt.Sprintf("(%c)", c)
		} else {
			s += fmt.Sprintf("%c", c)
		}
	}
	return s
}

func dump(cups []rune) string {
	g := &Game{cups: cups}
	return g.String()
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
