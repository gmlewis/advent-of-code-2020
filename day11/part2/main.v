// -*- compile-command: "v run main.v ../example1.txt ../input.txt"; -*-
import os

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
	grid   map[string]rune
	key_x  map[string]int
	key_y  map[string]int
}

fn process(filename string) {
	println('Processing $filename ...')
	mut puz := &Puzzle{}
	lines := os.read_lines(filename) or { panic(err) }
	for y, line in lines {
		puz.width = line.len
		puz.height++
		for x := 0; x < line.len; x++ {
			if line[x] == `L` {
				key := gen_key(x, y)
				puz.grid[key] = `L`
				puz.key_x[key] = x
				puz.key_y[key] = y
			}
		}
	}
	for puz.iterate() {}
	occupied := puz.occupied()
	println('Solution: $occupied')
}

fn gen_key(x int, y int) string {
	return '$x,$y'
}

fn (mut p Puzzle) iterate() bool {
	mut r := map[string]rune{}
	mut changes_made := false
	for k, v in p.grid {
		adj := p.count_adjacent_occupied(k)
		if v == `L` && adj == 0 {
			r[k] = `#`
			changes_made = true
		} else if v == `#` && adj >= 4 {
			r[k] = `L`
			changes_made = true
		} else {
			r[k] = v
		}
	}
	if changes_made {
		p.grid = r.clone()
	}
	return !changes_made
}

// As of 2021-09-05, V version 0.2.4, the following gives this error:
//
// $ v run main.v ../example1.txt ../input.txt
// main.v:93:6: error: undefined ident: `p`
//    91 |     t := fn (x int, y int) {
//    92 |         new_k := gen_key(x, y)
//    93 |         if p.grid[new_k] == `#` {
//       |            ^
//    94 |             adj++
//    95 |         }
// main.v:94:4: error: undefined ident: `adj`
//    92 |         new_k := gen_key(x, y)
//    93 |         if p.grid[new_k] == `#` {
//    94 |             adj++
//       |             ~~~
//    95 |         }
//    96 |     }
//
// Yet this worked fine when it was originally written.

fn (p &Puzzle) count_adjacent_occupied(k string) int {
	x := p.key_x[k]
	y := p.key_x[k]
	mut adj := 0

	t := fn (x int, y int) {
		new_k := gen_key(x, y)
		if p.grid[new_k] == `#` {
			adj++
		}
	}

	t(x - 1, y - 1)
	t(x, y - 1)
	t(x + 1, y - 1)
	t(x - 1, y)

	t(x + 1, y)
	t(x - 1, y + 1)
	t(x, y + 1)
	t(x + 1, y + 1)

	return adj
}

fn (p &Puzzle) occupied() int {
	mut count := 0
	for _, v in p.grid {
		if v == `#` {
			count++
		}
	}
	return count
}
