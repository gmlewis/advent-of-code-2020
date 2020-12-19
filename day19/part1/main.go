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
			rules[n] = &node{s: p1[1][2:3]}
			continue
		}

		rules[n] = parseRule(p1[1])
	}
	// log.Printf("got %v rules", len(rules))

	var sum int
	for _, msg := range strings.Split(msgs, "\n") {
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}

		n, ok := match(msg, 0, rules)
		logf("match rule 0 %v = (%v,%v)\n\n", msg, n, ok)
		if ok && n == len(msg) {
			sum++
		}
	}
	log.Printf("Solution: %v", sum)
}

func match(s string, ruleN int, rules map[int]*node) (int, bool) {
	logf("MATCH rule %v %v", ruleN, s)
	rule := rules[ruleN]

	if len(s) == 0 {
		return 0, true
	}

	if rule.s != "" {
		return 1, s[0:1] == rule.s
	}

	logf("Checking seq1 %v with %s", rule.seq1, s)
	n, ok := matchSeq(s, rule.seq1, rules)
	logf("matchSeq %v seq1 %v = (%v,%v)", s, rule.seq1, n, ok)
	if ok {
		return n, true
	}

	if len(rule.seq2) > 0 {
		logf("Checking seq2 %v with %s", rule.seq2, s)
		n, ok := matchSeq(s, rule.seq2, rules)
		logf("matchSeq %v seq2 %v = (%v,%v)", s, rule.seq2, n, ok)
		if ok {
			return n, true
		}
	}

	return 0, false
}

func matchSeq(s string, seq []int, rules map[int]*node) (int, bool) {
	if len(s) == 0 && len(seq) == 0 {
		logf("A: matchSeq - 0, true")
		return 0, true
	}
	if len(seq) == 0 {
		logf("A: matchSeq - 0, false (s=%v)", s)
		return 0, false
	}

	var numChars int
	for _, ruleN := range seq {
		n, ok := match(s, ruleN, rules)
		logf("sub match rule %v %v = (%v,%v)", ruleN, s, n, ok)
		if !ok {
			return numChars, false
		}
		numChars += n
		logf("s=%v len(s)=%v, numChars=%v", s, len(s), numChars)
		s = s[n:]
	}
	return numChars, true
}

func parseRule(s string) *node {
	rule := &node{}
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
