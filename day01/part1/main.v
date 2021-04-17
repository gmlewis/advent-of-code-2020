// -*- compile-command: "v run main.v ../input.txt"; -*-
import os
import strconv

fn main() {
	for arg in os.args[1..] {
		process(arg)
	}
	println('Done.')
}

fn process(filename string) {
	println('Processing $filename ...')
	lines := os.read_lines(filename) or { panic(err) }
	mut vals := map[string]bool{}
	for line in lines {
		v := strconv.atoi(line) or { panic(err) }
		if vals['$v'] {
			println('$v + ${2020 - v} = 2020\n$v * ${2020 - v} = ${v * (2020 - v)}')
			break
		}
		vals['${2020 - v}'] = true
	}
}
