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
	de    int
	dn    int
	epos  int
	npos  int
}

fn process(filename string) {
	println('Processing $filename ...')
	mut puz := &Puzzle{
		de: 1
		lines: os.read_lines(filename) or { panic(err) }
	}
	puz.iterate()
	dist := puz.manhattan()
	println('Solution: $dist')
}

fn (mut p Puzzle) iterate() {
	for line in p.lines {
		mut amt := strconv.atoi(line[1..]) or { panic(err) }
		match line[0] {
			`N` {
				p.npos += amt
			}
			`S` {
				p.npos -= amt
			}
			`E` {
				p.epos += amt
			}
			`W` {
				p.epos -= amt
			}
			`L` {
				for amt > 0 {
					p.de, p.dn = -p.dn, p.de
					amt -= 90
				}
			}
			`R` {
				for amt > 0 {
					p.de, p.dn = p.dn, -p.de
					amt -= 90
				}
			}
			`F` {
				p.npos += amt * p.dn
				p.epos += amt * p.de
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
