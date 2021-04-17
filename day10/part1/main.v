// -*- compile-command: "v run main.v ../example1.txt ../example2.txt ../input.txt"; -*-
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
	mut ints := []int{}
	for line in lines {
		n := strconv.atoi(line) or { panic(err) }
		ints << n
	}
	ints.sort()
	mut one_diffs := 1
	mut three_diffs := 1
	for i := 1; i < ints.len; i++ {
		match ints[i] - ints[i - 1] {
			1 { one_diffs++ }
			3 { three_diffs++ }
			else {}
		}
	}
	println('one_diffs=$one_diffs, three_diffs=$three_diffs, Solution: ${one_diffs * three_diffs}')
}
