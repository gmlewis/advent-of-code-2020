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
	spaces := os.read_lines(filename) or { panic(err) }
	mut ids := process_spaces(spaces)
	ids.sort()
	println('Highest seat id: ${ids[ids.len - 1]}')
}

fn process_spaces(lines []string) []int {
	mut result := []int{}
	for line in lines {
		v := space_id(line)
		result << v
	}
	return result
}

fn space_id(s string) int {
	mut v := s.replace('F', '0')
	v = v.replace('B', '1')
	v = v.replace('L', '0')
	v = v.replace('R', '1')
	return int(strconv.parse_int(v, 2, 64))
}
