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
		v := strconv.atoi(line)
		if find(v, vals) {
			break
		}
		vals['${2020 - v}'] = true
	}
}

fn find(v int, vals map[string]bool) bool {
	for d2, _ in vals {
		v2 := 2020 - strconv.atoi(d2)
		if vals['${v + v2}'] {
			println('$v + $v2 + ${2020 - v - v2} = 2020\n$v * $v2 * ${2020 - v - v2} = ${v * v2 *
				(2020 - v - v2)}')
			return true
		}
	}
	return false
}
