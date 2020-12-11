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
	mut ints := map[string]bool{}
	mut max := 0
	for line in lines {
		n := strconv.atoi(line)
		ints['$n'] = true
		if n > max {
			max = n
		}
	}
	mut cache := map[string]i64{}
	count := count_possibilities(0, ints, max, mut cache)
	println('Solution: $count')
}

fn count_possibilities(last_n int, ints map[string]bool, max int, mut cache map[string]i64) i64 {
	key := '$last_n'
	if key in cache {
		return cache[key]
	}
	if last_n == max {
		cache[key] = 1
		return 1
	}
	if last_n > max {
		cache[key] = 0
		return 0
	}
	mut result := i64(0)
	for i := 1; i <= 3; i++ {
		if '${last_n + i}' in ints {
			result += count_possibilities(last_n + i, ints, max, mut cache)
		}
	}
	cache[key] = result
	return result
}
