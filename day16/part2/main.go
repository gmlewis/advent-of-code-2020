// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "Verbose log messages")

	valsRE = regexp.MustCompile(`^(.*): (\d+)-(\d+) or (\d+)-(\d+)$`)
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		process(arg)
	}

	log.Printf("Done.")
}

func process(filename string) {
	logf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	groups := strings.Split(string(buf), "\n\n")
	puz := &Puzzle{fields: map[string]*valRanges{}, possible: map[int]map[string]bool{}, reject: map[int]map[string]bool{}, fieldName: map[int]string{}}
	puz.processRanges(groups[0])
	puz.errorRate(groups[2])
	puz.resolveFields()
	solution := puz.departureSolution(groups[1])

	log.Printf("Solution: %v", solution)
}

type Puzzle struct {
	fields    map[string]*valRanges
	fieldName map[int]string
	possible  map[int]map[string]bool
	reject    map[int]map[string]bool
}

type valRanges struct {
	fieldNum int
	notValid map[int]bool

	s1 int
	e1 int
	s2 int
	e2 int
}

func (p *Puzzle) resolveFields() {
	changesMade := true
	for changesMade {
		changesMade = false
		for k, v := range p.possible {
			logf("%v,%#v", k, v)
			if len(v) == 1 {
				var fieldName string
				for fn := range v {
					fieldName = fn
				}
				logf("%v is ABSOLUTELY fieldNum %v", fieldName, k)
				p.fieldName[k] = fieldName

				for k2, v2 := range p.possible {
					if k2 == k {
						continue
					}
					if v2[fieldName] {
						delete(p.possible[k2], fieldName)
						changesMade = true
					}
				}
			}
		}
	}
}

func (p *Puzzle) departureSolution(buf string) int64 {
	result := int64(1)
	line := strings.Split(buf, "\n")[1]
	for i, part := range strings.Split(line, ",") {
		fieldName, ok := p.fieldName[i]
		if !ok {
			continue
		}
		if strings.Contains(fieldName, "departure") {
			v, err := strconv.ParseInt(part, 10, 64)
			check("v2: %v", err)
			result *= v
			logf("%v,%v: %q, result=%v", i, v, fieldName, result)
		}
	}
	return result
}

func (p *Puzzle) processRanges(buf string) {
	for _, line := range strings.Split(buf, "\n") {
		m := valsRE.FindStringSubmatch(line)
		if len(m) != 6 {
			log.Fatalf("bad input/regexp: %v", line)
		}
		logf("line: %v", line)
		s1, err := strconv.Atoi(m[2])
		check("s1: %v", err)
		e1, err := strconv.Atoi(m[3])
		check("e1: %v", err)
		s2, err := strconv.Atoi(m[4])
		check("s2: %v", err)
		e2, err := strconv.Atoi(m[5])
		check("e2: %v", err)
		p.fields[m[1]] = &valRanges{s1: s1, e1: e1, s2: s2, e2: e2, notValid: map[int]bool{}}
		logf("fields[%q] = %#v", m[1], p.fields[m[1]])
	}
}

func (p *Puzzle) errorRate(buf string) int {
	var result int
	for _, line := range strings.Split(buf, "\n")[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		for i, part := range parts {
			v, err := strconv.Atoi(part)
			check("v: %v - %v", line, err)
			if !p.isValid(v, i) {
				result += v
			}
		}
	}
	return result
}

func (p *Puzzle) isValid(n, fieldNum int) bool {
	var valid bool
	for _, v := range p.fields {
		if (n >= v.s1 && n <= v.e1) || (n >= v.s2 && n <= v.e2) {
			valid = true
			break
		}
	}
	if !valid {
		return false
	}

	for k, v := range p.fields {
		if (n >= v.s1 && n <= v.e1) || (n >= v.s2 && n <= v.e2) {
			if p.reject[fieldNum][k] {
				continue
			}

			v.fieldNum = fieldNum
			if _, ok := p.possible[fieldNum]; !ok {
				p.possible[fieldNum] = map[string]bool{}
			}

			//			if strings.Contains(k, "departure") {
			p.possible[fieldNum][k] = true
			logf("%v is possibly field %v due to (%v,%v)", k, fieldNum, n, fieldNum)
			//			}
			valid = true
			continue
		}
		if _, ok := p.reject[fieldNum]; !ok {
			p.reject[fieldNum] = map[string]bool{}
		}
		p.reject[fieldNum][k] = true
		if _, ok := p.possible[fieldNum]; !ok {
			p.possible[fieldNum] = map[string]bool{}
		}
		delete(p.possible[fieldNum], k)
		logf("%v is definitely not field %v due to (%v,%v)", k, fieldNum, n, fieldNum)
	}
	return valid
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
