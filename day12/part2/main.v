// -*- compile-command: "v run main.v ../example1.txt ../input.txt"; -*-
import os
import strconv

fn main() {
	for arg in os.args[1..] {
		process(arg)
	}
	println('Done.')
}

struct Puzzle {
	lines []string
mut:
	wepos int
	wnpos int
	epos  int
	npos  int
}

fn process(filename string) {
	println('Processing $filename ...')
	mut puz := &Puzzle{
		wepos: 10
		wnpos: 1
		lines: os.read_lines(filename) or { panic(err) }
	}
	puz.iterate()
	dist := puz.manhattan()
	println('Solution: $dist')
}

fn (mut p Puzzle) iterate() {
	for line in p.lines {
		mut amt := strconv.atoi(line[1..])
		match line[0] {
			`N` {
				p.wnpos += amt
			}
			`S` {
				p.wnpos -= amt
			}
			`E` {
				p.wepos += amt
			}
			`W` {
				p.wepos -= amt
			}
			`L` {
				for amt > 0 {
					p.wepos, p.wnpos = -p.wnpos, p.wepos
					amt -= 90
				}
			}
			`R` {
				for amt > 0 {
					p.wepos, p.wnpos = p.wnpos, -p.wepos
					amt -= 90
				}
			}
			`F` {
				p.npos += amt * p.wnpos
				p.epos += amt * p.wepos
			}
			else {}
		}
	}
}

fn (p &Puzzle) manhattan() int {
	x := if p.epos >= 0 { p.epos } else { -p.epos }
	y := if p.npos >= 0 { p.npos } else { -p.npos }
	return x + y
}
