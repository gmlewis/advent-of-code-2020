// -*- compile-command: "v run main.v ../example1.txt ../input.txt"; -*-
import os

const (
	down  = 1
	right = 3
)

fn main() {
	for arg in os.args[1..] {
		process(arg)
	}
	println('Done.')
}

struct Puzzle {
mut:
	width  int
	height int
	grid   map[string]bool
}

fn process(filename string) {
	println('Processing $filename ...')
	mut puz := &Puzzle{}
	lines := os.read_lines(filename) or { panic(err) }
	for y, line in lines {
		puz.width = line.len
		puz.height++
		for x := 0; x < line.len; x++ {
			if line[x] == `#` {
				key := gen_key(x, y)
				puz.grid[key] = true
			}
		}
	}
	mut pos_x, mut pos_y, mut count := 0, 0, 0
	for y := 0; y < puz.height; y++ {
		if puz.lookup(pos_x, pos_y) {
			count++
		}
		pos_x += right
		pos_y += down
	}
	println('Found $count trees')
}

fn gen_key(x int, y int) string {
	return '$x,$y'
}

fn (p &Puzzle) lookup(x int, y int) bool {
	key := gen_key(x % p.width, y)
	return p.grid[key]
}
