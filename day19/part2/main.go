// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
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

	parts := strings.Split(string(buf), "\n\n")
	msgs := parts[1]

	lines := strings.Split(parts[0], "\n")
	rules := map[int]*node{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		p1 := strings.Split(line, ":")
		n, err := strconv.Atoi(p1[0])
		check("n: %v", err)
		if strings.HasPrefix(p1[1], ` "`) {
			s := p1[1][2:3]
			rules[n] = &node{s: s, matches: map[string]bool{s: true}}
			// log.Printf("rules[%v] = %#v", n, rules[n])
			continue
		}

		if n == 8 {
			p1[1] = " 42 | 42 8"
		} else if n == 11 {
			p1[1] = " 42 31 | 42 11 31"
		}

		rules[n] = parseRule(p1[1])
		// log.Printf("rules[%v] = %#v", n, rules[n])
	}
	// log.Printf("got %v rules", len(rules))

	var sum int
	for i, msg := range strings.Split(msgs, "\n") {
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}

		if n, ok := match(msg, []int{0}, rules); ok && n == len(msg) {
			sum++
			log.Printf("line %v: MATCH #%v! %v", i, sum, msg)
		}
	}
	log.Printf("Solution: %v", sum)
}

func match(s string, seq []int, rules map[int]*node) (int, bool) {
	log.Printf("ENTER: match(%q[%v], seq=%v)", s, len(s), seq)
	if len(s) == 0 {
		log.Printf("A: LEAVE: 0, %v", len(seq) == 0)
		return 0, len(seq) == 0
	}
	if len(seq) == 0 {
		log.Printf("B: LEAVE: 0, false")
		return 0, false
	}

	rule := rules[seq[0]]
	if len(s) == 1 && len(seq) == 1 && rule.s == s {
		log.Printf("C1: LEAVE: 1, true")
		return 1, true
	}

	if len(seq) == 1 && rule.s == s[0:1] {
		log.Printf("C2: LEAVE: 1, true")
		return 1, true
	}

	tmpseq := rule.seq1[:]
	tmps := s
	for len(tmpseq) > 0 {
		log.Printf("CALLING: match(%q[%v], %v) ...", tmps, len(tmps), tmpseq)
		n, ok := match(tmps, tmpseq, rules)
		log.Printf("GOT: match(%q[%v], %v) = (%v,%v)", tmps, len(tmps), tmpseq, n, ok)
		if !ok {
			break
		}
		tmps = tmps[n:]
		tmpseq = tmpseq[1:]
		log.Printf("tmps=%q[%v], tmpseq=%v", tmps, len(tmps), tmpseq)
		if len(tmpseq) == 0 {
			log.Printf("D: LEAVE: %v, true", len(s)-len(tmps))
			return len(s) - len(tmps), true
		}
	}

	tmpseq = rule.seq2[:]
	tmps = s
	for len(tmpseq) > 0 {
		log.Printf("CALLING: match(%q[%v], %v) ...", tmps, len(tmps), tmpseq)
		n, ok := match(tmps, tmpseq, rules)
		log.Printf("GOT: match(%q[%v], %v) = (%v,%v)", tmps, len(tmps), tmpseq, n, ok)
		if !ok {
			break
		}
		tmps = tmps[n:]
		tmpseq = tmpseq[1:]
		log.Printf("tmps=%q[%v], tmpseq=%v", tmps, len(tmps), tmpseq)
		if len(tmpseq) == 0 {
			log.Printf("D: LEAVE: %v, true", len(s)-len(tmps))
			return len(s) - len(tmps), true
		}
	}

	log.Printf("F: LEAVE: match(%q[%v], seq=%v), rule=%v: 0, false", s, len(s), seq, rule)

	return 0, false
}

func parseRule(s string) *node {
	rule := &node{matches: map[string]bool{}}
	seq := &rule.seq1
	parts := strings.Split(strings.TrimSpace(s), " ")
	for _, part := range parts {
		if part == "|" {
			seq = &rule.seq2
			continue
		}
		n, err := strconv.Atoi(part)
		check("n: %v", err)
		*seq = append(*seq, n)
	}
	// logf("parseRule: %v = %#v", s, *rule)
	return rule
}

type node struct {
	s string

	matches map[string]bool

	seq1 []int
	seq2 []int
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
