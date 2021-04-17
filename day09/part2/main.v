// -*- compile-command: "v run main.v -p 5 ../example1.txt && v run main.v -p 25 ../input.txt"; -*-
import os
import strconv

fn main() {
	if os.args.len != 4 || os.args[1] != '-p' {
		println('usage: v run main.v -p 25 ../input.txt')
		return
	}
	preamble := strconv.atoi(os.args[2]) or { panic(err) }
	process(os.args[3], preamble)
	println('Done.')
}

fn process(filename string, preamble int) {
	println('Processing $filename ...')
	lines := os.read_lines(filename) or { panic(err) }
	mut ints := []int{}
	for line in lines {
		n := strconv.atoi(line) or { panic(err) }
		ints << n
	}
	for i := preamble; i < ints.len; i++ {
		if found_solution(i, ints, preamble) {
			break
		}
	}
}

fn found_solution(n int, ints []int, preamble int) bool {
	v := ints[n]
	mut seen := map[string]bool{}
	for i := n - preamble; i < n; i++ {
		d := v - ints[i]
		if seen['$d'] {
			return false
		}
		seen['${ints[i]}'] = true
	}
	println('Part 1 Solution: $v')
	for i := n - 1; i >= preamble - 2; i-- {
		mut sum := 0
		mut min := ints[i]
		mut max := ints[i]
		for j := 0; j < preamble - 1; j++ {
			sample := ints[i - j]
			if sample < min {
				min = sample
			}
			if sample > max {
				max = sample
			}
			sum += sample
			if sum > v {
				break
			}
			if sum == v {
				println('Part 2 solution: $min + $max = ${min + max}')
				return true
			}
		}
	}
	return true
}
