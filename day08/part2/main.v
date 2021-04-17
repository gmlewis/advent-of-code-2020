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
	mut program := []&Instruction{}
	for line in lines {
		parts := line.split(' ')
		ins := &Instruction{
			op: parts[0]
			arg: strconv.atoi(parts[1]) or { panic(err) }
		}
		program << ins
	}
	cpu.initialize(program)
	recording, _ := cpu.execute()
	for i := recording.len - 1; i >= 0; i-- {
		cpu.initialize(program)
		ip := recording[i]
		ins := program[ip]
		match ins.op {
			'nop' {
				cpu.program[ip] = &Instruction{
					op: 'jmp'
					arg: ins.arg
				}
			}
			'jmp' {
				cpu.program[ip] = &Instruction{
					op: 'nop'
					arg: ins.arg
				}
			}
			else {
				continue
			}
		}
		_, done := cpu.execute()
		if done {
			break
		}
	}
}

fn (mut c CPU) initialize(program []&Instruction) {
	c.accumulator = 0
	c.ip = 0
	c.program = program[0..program.len]
}

fn (mut c CPU) execute() ([]int, bool) {
	mut seen := map[string]bool{}
	mut recorder := []int{}
	for {
		if c.ip >= c.program.len {
			println('Normal termination: Solution Accumulator = $c.accumulator')
			return recorder, true
		}
		if seen['$c.ip'] {
			// println('Infinite loop: Accumulator = $c.accumulator')
			return recorder, false
		}
		seen['$c.ip'] = true
		ins := c.program[c.ip]
		recorder << c.ip
		op := ins.op
		arg := ins.arg
		match op {
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
	return recorder, false
}
