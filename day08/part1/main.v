// -*- compile-command: "v run main.v ../example1.txt ../input.txt"; -*-
import os
import strconv

fn main() {
	for arg in os.args[1..] {
		process(arg)
	}
	println('Done.')
}

struct CPU {
mut:
	accumulator int
	ip          int
	program     []&Instruction
}

struct Instruction {
	op  string
	arg int
}

fn process(filename string) {
	println('Processing $filename ...')
	lines := os.read_lines(filename) or { panic(err) }
	mut cpu := &CPU{}
	for line in lines {
		parts := line.split(' ')
		ins := &Instruction{
			op: parts[0]
			arg: strconv.atoi(parts[1]) or { panic(err) }
		}
		cpu.program << ins
	}
	cpu.execute()
}

fn (mut c CPU) execute() {
	mut seen := map[string]bool{}
	for {
		if seen['$c.ip'] {
			println('Solution Accumulator = $c.accumulator')
			break
		}
		seen['$c.ip'] = true
		arg := c.program[c.ip].arg
		match c.program[c.ip].op {
			'nop' {
				c.ip++
			}
			'acc' {
				c.accumulator += arg
				c.ip++
			}
			'jmp' {
				c.ip += arg
			}
			else {}
		}
	}
}
