// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "Verbose log messages")

	foodRE = regexp.MustCompile(`^(.*?) \(contains (.*?)\)$`)
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
	var foods []*foodT
	foodsByAlrgn := map[string]*alrgnT{}
	allIngredCounts := map[string]int{}
	ingredIsAllergen := map[string]string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		food := parseFood(line)
		for ingred := range food.ingred {
			allIngredCounts[ingred]++
		}
		for alrgn := range food.alrgns {
			if _, ok := foodsByAlrgn[alrgn]; !ok {
				foodsByAlrgn[alrgn] = &alrgnT{foods: map[int]bool{}}
			}
			foodsByAlrgn[alrgn].foods[len(foods)] = true
		}
		foods = append(foods, food)
	}

	changesMade := true
	for changesMade {
		changesMade = false
		for k, v := range foodsByAlrgn {
			if v.ingred != "" {
				continue
			}
			logf("BEFORE foodsByAlrgn[%v]: %v", k, v)
			ingred := v.findIngredWithAlrgn(foods, ingredIsAllergen)
			logf("AFTER foodsByAlrgn[%v]: %v", k, v)
			if ingred != "" {
				v.ingred = ingred
				ingredIsAllergen[ingred] = k
				changesMade = true
			}
		}
	}

	var allergens []string
	for ingred := range allIngredCounts {
		if ingredIsAllergen[ingred] != "" {
			allergens = append(allergens, ingred)
		}
	}
	sort.Slice(allergens, func(a, b int) bool { return ingredIsAllergen[allergens[a]] < ingredIsAllergen[allergens[b]] })
	log.Printf("Solution: %v", fmt.Sprintf("%v", strings.Join(allergens, ",")))
}

type alrgnT struct {
	ingred string
	foods  map[int]bool
}

func (a *alrgnT) findIngredWithAlrgn(foods []*foodT, ingredIsAllergen map[string]string) string {
	ingCount := map[string]int{}
	for k := range a.foods {
		for ingred := range foods[k].ingred {
			ingCount[ingred]++
		}
	}

	logf("a: %v", *a)
	possibilities := map[string]bool{}
	var lastPossibility string
	for k, v := range ingCount {
		if ingredIsAllergen[k] != "" {
			logf("Ingredient %v (%v) is already identified as an allergen. Skipping.", k, ingredIsAllergen[k])
			continue
		}
		if v == len(a.foods) {
			logf("Found (%v,%v)", k, v)
			possibilities[k] = true
			lastPossibility = k
		}
	}

	if len(possibilities) == 1 {
		logf("%v is ABSOLUTELY the allergen.", lastPossibility)
		return lastPossibility
	}

	// log.Fatalf("findIngredWithAlrgn: %v", a)
	return ""
}

type foodT struct {
	ingred map[string]bool
	alrgns map[string]bool
}

func parseFood(s string) *foodT {
	food := &foodT{ingred: map[string]bool{}, alrgns: map[string]bool{}}
	m := foodRE.FindStringSubmatch(s)
	if len(m) != 3 {
		log.Fatalf("bad foodRE: %v", s)
	}

	for _, ingred := range strings.Split(m[1], " ") {
		food.ingred[ingred] = true
	}
	for _, alrgn := range strings.Split(m[2], ", ") {
		food.alrgns[alrgn] = true
	}
	return food
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
