// -*- compile-command: "v run main.v ../input.txt"; -*-
import os

fn main() {
	for arg in os.args[1..] {
		process(arg)
	}
	println('Done.')
}

fn process(filename string) {
	println('Processing $filename ...')
	buf := os.read_file(filename) or { panic(err) }
	groups := buf.split('\n\n')
	mut count := 0
	for group in groups {
		v := group.replace('\n', '')
		mut runes := map[string]int{}
		for i := 0; i < v.len; i++ {
			r := v[i..i + 1]
			runes[r]++
		}
		count += runes.len
	}
	println('Solution: $count')
}
