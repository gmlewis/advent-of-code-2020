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
		lines := group.trim_space().split('\n')
		mut runes := map[string]int{}
		for line in lines {
			for i := 0; i < line.len; i++ {
				r := line[i..i + 1]
				runes[r]++
				if runes[r] == lines.len {
					count++
				}
			}
		}
	}
	println('Solution: $count')
}
