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

	lineRE     = regexp.MustCompile(`^(.*?) bags contain (.*)\.$`)
	containsRE = regexp.MustCompile(`,?\s*(\d+) (.*?) bags?`)
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
	var program []*Instruction

	lines := readLines(filename)
	logf("lines=%v", lines)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		op := parts[0]
		arg, err := strconv.Atoi(parts[1])
		check("strconv.Atoi: %v", err)
		program = append(program, &Instruction{op: op, arg: arg})
	}

	cpu.initialize(program)
	recording, _ := cpu.execute()

	for i := len(recording) - 1; i >= 0; i-- {
		cpu.initialize(program)
		ip := recording[i]
		ins := program[ip]
		switch ins.op {
		case "nop":
			log.Printf("Experiment: Changing op at ip=%v from nop to jmp", ip)
			cpu.program[ip].op = "jmp"
		case "jmp":
			log.Printf("Experiment: Changing op at ip=%v from jmp to nop", ip)
			cpu.program[ip].op = "nop"
		default:
			continue
		}

		if _, done := cpu.execute(); done {
			break
		}
	}
}

func (cpu *CPU) initialize(program []*Instruction) {
	cpu.accumulator = 0
	cpu.ip = 0
	cpu.program = program[:]
}

func (cpu *CPU) execute() ([]int, bool) {
	seen := map[int]bool{}
	var recorder []int

	for {
		if cpu.ip >= len(cpu.program) {
			log.Printf("Normal termination: Solution Accumulator = %v", cpu.accumulator)
			return recorder, true
		}
		if seen[cpu.ip] {
			log.Printf("Infinite loop: Accumulator = %v", cpu.accumulator)
			return recorder, false
		}
		seen[cpu.ip] = true

		instruction := cpu.program[cpu.ip]
		recorder = append(recorder, cpu.ip)
		op := instruction.op
		arg := instruction.arg
		logf("%3d: %v %-3d |", cpu.ip, op, arg)
		switch op {
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
