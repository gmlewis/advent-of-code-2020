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

	lines := strings.Split(string(buf), "\n")
	var sum int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v := eval(line)
		sum += v
		logf("%v = %v, sum = %v", line, v, sum)
	}

	log.Printf("Solution: %v", sum)
}

type node struct {
	s string

	isNum bool
	val   int
}

func (n *node) String() string {
	if n.isNum {
		return fmt.Sprintf("%v", n.val)
	}
	return n.s
}

func eval(s string) int {
	var stack []*node
	for _, v := range s {
		if v == ' ' {
			continue
		}
		n := &node{s: fmt.Sprintf("%c", v)}
		if v >= '0' && v <= '9' {
			n.isNum = true
			n.val = int(v - '0')
		}
		stack = append(stack, n)
	}

	result, _ := processParens(stack)
	v := processExpr(result)
	return v
}

func processParens(stack []*node) ([]*node, int) {
	// log.Printf("processParens: %v", stack)
	var result []*node
	for i := 0; i < len(stack); {
		n := stack[i]
		// log.Printf("i=%v: %v", i, n)
		if n.isNum || n.s == "+" || n.s == "*" {
			result = append(result, n)
			i++
			continue
		}
		if n.s == "(" {
			sub, index := processParens(stack[i+1:])
			v := processExpr(sub)
			result = append(result, &node{isNum: true, val: v})
			i += index + 2
			continue
		}
		if n.s == ")" {
			return result, i
		}
	}
	return result, len(stack)
}

func processExpr(stack []*node) int {
	// log.Printf("process: %v", stack)
	var sp int
	var result int
	for sp < len(stack) {
		if stack[sp].isNum {
			sp++
			continue
		}
		if stack[sp].s == "+" {
			result = stack[sp-1].val + stack[sp+1].val
			stack[sp+1].val = result
			sp = sp + 1
			continue
		}
		if stack[sp].s == "*" {
			result = stack[sp-1].val * stack[sp+1].val
			stack[sp+1].val = result
			sp = sp + 1
		}
	}
	return result
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
