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

	for i := len(try) - 1; i >= 0; i-- {
		cpu.execute(try[i])
	}

	// fmt.Printf("Solution: %v\n", count-1)
}

func (cpu *CPU) execute(change int) {
	cpu.accumulator = 0
	cpu.ip = 0
	seen := map[int]bool{}

	log.Printf("Program length = %v", len(cpu.program))
	if len(cpu.program) > 491 {
		// Experiment: 491 nop => jmp
		// cpu.program[491].op = "jmp"
		// Experiment: 492 jmp => nop
		// cpu.program[492].op = "nop"
		// Experiment: 488 jmp => nop
		// cpu.program[488].op = "nop"
		// Experiment: 486 nop => jmp
		// cpu.program[486].op = "jmp"
		// Experiment: 151 jmp => nop
		// cpu.program[151].op = "nop"
		// Experiment: 609 jmp => nop
		// cpu.program[609].op = "nop"
		// Experiment: 607 nop => jmp
		// cpu.program[607].op = "jmp"
		// Experiment: 384 jmp => nop
		// cpu.program[384].op = "nop"
		if cpu.program[change].op == "nop" {
			log.Printf("Experiment: Changing op at ip=%v from nop to jmp", change)
			cpu.program[change].op = "jmp"
		} else {
			log.Printf("Experiment: Changing op at ip=%v from jmp to nop", change)
			cpu.program[change].op = "nop"
		}
	}

	for {
		if cpu.ip >= len(cpu.program) {
			log.Fatalf("Normal termination: Solution Accumulator = %v", cpu.accumulator)
			break
		}
		if seen[cpu.ip] {
			log.Printf("Solution Accumulator = %v", cpu.accumulator)
			break
		}
		seen[cpu.ip] = true

		op := cpu.program[cpu.ip].op
		arg := cpu.program[cpu.ip].arg
		log.Printf("%3d: %v %-3d |", cpu.ip, op, arg)
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

var try = []int{
	3,   // : jmp 292 |
	296, // : nop 31  |
	299, // : jmp -208 |
	92,  // : jmp 146 |
	238, // : nop -33 |
	239, // : nop -167 |
	240, // : jmp -174 |
	66,  // : nop 363 |
	70,  // : jmp 379 |
	451, // : nop -227 |
	452, // : jmp 116 |
	568, // : jmp -265 |
	303, // : nop -170 |
	307, // : jmp 295 |
	602, // : jmp -574 |
	28,  // : jmp 56  |
	84,  // : jmp 107 |
	193, // : nop 409 |
	195, // : jmp 83  |
	281, // : jmp 64  |
	346, // : jmp -114 |
	232, // : jmp -216 |
	18,  // : jmp 162 |
	180, // : jmp 106 |
	289, // : jmp 231 |
	520, // : jmp 98  |
	618, // : nop -164 |
	619, // : jmp -462 |
	159, // : jmp 110 |
	271, // : nop 32  |
	272, // : jmp 261 |
	535, // : jmp -439 |
	96,  // : jmp 482 |
	581, // : jmp -376 |
	205, // : nop -20 |
	206, // : jmp -152 |
	54,  // : jmp 1   |
	55,  // : jmp 484 |
	540, // : nop -254 |
	542, // : jmp -72 |
	470, // : jmp -448 |
	23,  // : jmp 385 |
	409, // : jmp 213 |
	623, // : jmp -438 |
	185, // : nop 354 |
	188, // : jmp 63  |
	251, // : nop 76  |
	254, // : jmp -150 |
	105, // : nop -33 |
	108, // : jmp 157 |
	267, // : jmp 345 |
	612, // : jmp -489 |
	127, // : jmp 81  |
	208, // : nop -161 |
	209, // : nop -76 |
	211, // : nop -202 |
	212, // : jmp -98 |
	114, // : jmp 219 |
	333, // : nop -254 |
	335, // : nop 277 |
	336, // : jmp 105 |
	442, // : jmp -40 |
	403, // : nop -125 |
	404, // : nop -7  |
	406, // : jmp -266 |
	140, // : jmp 1   |
	144, // : jmp 301 |
	447, // : jmp -77 |
	371, // : jmp -12 |
	359, // : jmp -241 |
	118, // : jmp 1   |
	119, // : jmp 278 |
	398, // : jmp -174 |
	227, // : jmp -187 |
	40,  // : nop 465 |
	42,  // : jmp 555 |
	599, // : jmp -258 |
	342, // : jmp -282 |
	60,  // : jmp 258 |
	322, // : jmp 225 |
	550, // : jmp -305 |
	248, // : jmp 323 |
	572, // : nop -340 |
	573, // : jmp -62 |
	515, // : jmp -483 |
	33,  // : jmp 223 |
	257, // : jmp 197 |
	455, // : jmp -74 |
	382, // : nop 76  |
	384, // : jmp 222 |
	607, // : nop -493 |
	609, // : jmp -459 |
	151, // : jmp 333 |
	486, // : nop -241 |
	488, // : jmp 2   |
	491, // : nop -351 |
	492, // : jmp -254 |
}
