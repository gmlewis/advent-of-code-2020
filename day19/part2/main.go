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

		n, lm, ok := match(msg, 0, -1, rules)
		logf("match rule 0 %v = (%v,lastMatch=%v,%v)\n\n", msg, n, lm, ok)
		if ok && n == len(msg) {
			sum++
			log.Printf("line %v: MATCH #%v! %v", i, sum, msg)
		}
	}
	log.Printf("Solution: %v", sum)
}

func match(s string, ruleN int, lastMatch int, rules map[int]*node) (int, int, bool) {
	rule := rules[ruleN]
	logf("ENTER: MATCH ruleN=%v lastMatch=%v, %v[%v]: %v", ruleN, lastMatch, s, len(s), *rule)

	cacheMatch := -1
	for k := range rule.matches {
		if len(s) >= len(k) && strings.HasPrefix(s, k) {
			logf("CACHE: LEAVE: MATCH ruleN=%v lastMatch=%v, %v[%v]: %v = (%v,true)", ruleN, lastMatch, s, len(s), *rule, len(k))
			cacheMatch = len(k)
			break
		}
	}
	if cacheMatch >= 0 {
		return cacheMatch, ruleN, true
	}

	if len(s) == 0 {
		logf("A: LEAVE: MATCH ruleN=%v lastMatch=%v, %v[%v]: %v = (0,false)", ruleN, lastMatch, s, len(s), *rule)
		return 0, lastMatch, false
	}

	if rule.s != "" {
		logf("B: LEAVE: MATCH ruleN=%v lastMatch=%v, %v[%v]: %v = (1,%v)", ruleN, lastMatch, s, len(s), *rule, s[0:1] == rule.s)
		if s[0:1] == rule.s {
			lastMatch = ruleN
		}
		return 1, lastMatch, s[0:1] == rule.s
	}

	logf("Checking seq1 %v with %s", rule.seq1, s)
	n, lm, ok := matchSeq(s, rule.seq1, lastMatch, rules)
	logf("matchSeq %v seq1 %v = (%v,lastMatch=%v,%v)", s, rule.seq1, n, lm, ok)
	if ok {
		lastMatch = lm
		rule.matches[s[:n]] = true
		if ruleN == 8 || ruleN == 11 {
			logf("C: ENTER SPECIAL CASE!!! ruleN=%v s=%v[%v]: %v = (%v,lastMatch=%v,true)", ruleN, s, len(s), *rule, n, lastMatch)
			for len(s) > n {
				tmp := s[n:]
				logf("ENTER ITERATE SPECIAL CASE!!! calling match(%v[%v],%v,%v): %v", tmp, len(tmp), 42, lastMatch, *rule)
				if nv, lm, ok := match(tmp, 42, lastMatch, rules); ok {
					logf("LEAVE ITERATE SPECIAL CASE!!! %v %v[%v]: %v = (%v,lastMatch=%v,true)", 42, tmp, len(tmp), *rule, nv, lm)
					lastMatch = lm
					n += nv
				} else {
					break
				}
			}
			logf("C: LEAVE SPECIAL CASE!!! ruleN=%v lastMatch=%v, %v[%v]: %v = (%v,lastMatch=%v,true)", ruleN, lastMatch, s, len(s), *rule, n, lastMatch)
		} else {
			logf("C: LEAVE: MATCH ruleN=%v lastMatch=%v, %v[%v]: %v = (%v,lastMatch=%v,true)", ruleN, lastMatch, s, len(s), *rule, n, lastMatch)
		}

		return n, lastMatch, true
	} else {
		if ruleN == 11 && (lastMatch == 8 || lastMatch == 42) {
			logf("C1: ENTER SUPER SPECIAL CASE!!! %v %v: %v = (%v,lastMatch=%v,true)", ruleN, s, *rule, n, lastMatch)
			logf("rules[8].matches = %v", rules[8].matches)
			logf("rules[11].matches = %v", rules[11].matches)
			logf("rules[42].matches = %v", rules[42].matches)
			for k := range rules[42].matches {
				if strings.HasPrefix(s, k) {
					log.Printf("rules[42].matches[%v] !!!", k)
					s = s[len(k):]
					n = len(k)
					break
				}
			}
			if nv, lm, ok := match(s, 31, lastMatch, rules); ok {
				logf("C1: WOOHOO SPECIAL CASE!!! %v %v: %v = (%v,lastMatch=%v,true)", 31, s, *rule, n+nv, lm)
				lastMatch = lm
				return n + nv, 31, true
			}
			logf("C1: LEAVE SUPER SPECIAL CASE!!! %v %v: %v = (%v,lastMatch=%v,true)", ruleN, s, *rule, n, lastMatch)
		}
	}

	if len(rule.seq2) > 0 {
		logf("Checking seq2 %v with %s", rule.seq2, s)
		n, lm, ok := matchSeq(s, rule.seq2, lastMatch, rules)
		logf("matchSeq %v seq2 %v = (%v,lastMatch=%v,%v)", s, rule.seq2, n, lm, ok)
		if ok {
			lastMatch = lm
			logf("D: LEAVE: MATCH rule %v %v: %v = (%v,lastMatch=%v,true)", ruleN, s, *rule, n, lm)
			rule.matches[s[:n]] = true
			return n, lastMatch, true
		}
	}

	logf("E0: LEAVE: MATCH rule %v %v: %v = (0,lastMatch=%v,false)", ruleN, s, *rule, lastMatch)
	return 0, lastMatch, false
}

func matchSeq(s string, seq []int, lastMatch int, rules map[int]*node) (int, int, bool) {
	if len(s) == 0 && len(seq) == 0 {
		logf("A: matchSeq - 0, true")
		return 0, lastMatch, true
	}
	if len(seq) == 0 {
		logf("A: matchSeq - 0, false (s=%v)", s)
		return 0, lastMatch, false
	}

	var numChars int
	for _, ruleN := range seq {
		n, lm, ok := match(s, ruleN, lastMatch, rules)
		if !ok {
			logf("sub match rule %v %v = (%v,lastMatch=%v,%v)", ruleN, s, n, lm, ok)
			return 0, lastMatch, false
		}
		logf("sub match rule %v %v = (%v,lastMatch=%v...but setting to %v,%v)", ruleN, s, n, lm, ruleN, ok)
		lastMatch = ruleN
		numChars += n
		logf("s=%v[%v], numChars=%v", s, len(s), numChars)
		s = s[n:]
	}
	return numChars, lastMatch, true
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
