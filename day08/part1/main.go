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

type CPU struct {
	accumulator int
	ip          int
	program     []*Instruction
}

type Instruction struct {
	op  string
	arg int
}

func process(filename string) {
	cpu := &CPU{}

	lines := readLines(filename)
	// log.Printf("lines=%v", lines)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		op := parts[0]
		arg, err := strconv.Atoi(parts[1])
		check("strconv.Atoi: %v", err)
		cpu.program = append(cpu.program, &Instruction{op: op, arg: arg})
	}

	cpu.execute()

	// fmt.Printf("Solution: %v\n", count-1)
}

func (cpu *CPU) execute() {
	seen := map[int]bool{}

	for {
		if seen[cpu.ip] {
			log.Printf("Solution Accumulator = %v", cpu.accumulator)
			break
		}
		seen[cpu.ip] = true

		arg := cpu.program[cpu.ip].arg
		switch cpu.program[cpu.ip].op {
		case "nop":
			cpu.ip++
		case "acc":
			cpu.accumulator += arg
			cpu.ip++
		case "jmp":
			cpu.ip += arg
		}
	}
}

func readLines(filename string) []string {
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)
	lines := strings.Split(string(buf), "\n")
	return lines
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
