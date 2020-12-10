// -*- compile-command: "v run main.v ../input.txt"; -*-
import os
import regex
import strconv

fn main() {
	for arg in os.args[1..] {
		process(arg)
	}
	println('Done.')
}

struct Contains {
	quant int
	color string
}

fn process(filename string) {
	println('Processing $filename ...')
	lines := os.read_lines(filename) or { panic(err) }
	mut line_re := regex.new()
	// NOTE: As of 2020-12-09, the regex support is too buggy to use.
	line_re.compile_opt(r'^(.*?) bags contain (.*)\.$')
	mut contains_re := regex.new()
	contains_re.compile_opt(r',?\s*(\d+) (.*?) bags?')
	mut rules := map[string][]&Contains{}
	for line in lines {
		line_re.match_string(line)
		m1 := re_field(line_re, line, 0)
		m2 := re_field(line_re, line, 1)
		// println('$line: m1=$m1, m2=$m2')
		if m2 == 'no other bags' {
			rules[re_field(line_re, line, 0)] = []
			continue
		}
		m := contains_re.find_all(m2)
		println('line=$line')
		println('m1=$m1')
		println('m2=$m2')
		println('m=$m')
		for i := 0; i < m.len; i += 4 {
			f0 := m2[m[i]..m[i + 1]]
			f1 := m2[m[i + 1]..m[i + 2]]
			q := strconv.atoi(f0)
			println('$m2: $f0, $f1, $q')
			rules[m1] << &Contains{
				quant: q
				color: f1
			}
		}
	}
	mut count := 0
	for color, _ in rules {
		mut seen := map[string]bool{}
		if can_contain(color, 'shiny gold', rules, mut seen) {
			count++
		}
	}
	println('Solution: $count')
}

fn can_contain(key string, test string, rules map[string][]&Contains, mut seen map[string]bool) bool {
	println('can_contain($key, $test, $rules, $seen)')
	for v in rules[key] {
		if v.color == test {
			return true
		}
		if seen[v.color] {
			continue
		}
		seen[v.color] = true
		if can_contain(v.color, test, rules, mut seen) {
			return true
		}
	}
	return false
}

fn re_field(re regex.RE, line string, n int) string {
	gp := re.get_group_list()
	// println('${gp[n]}, line=$line, len=$line.len')
	start := gp[n].start
	// regex bug workaround - https://github.com/vlang/v/issues/7227
	mut end := gp[n].end
	if end > line.len {
		end = line.len
	}
	return line[start..end]
}
