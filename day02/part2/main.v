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

fn process(filename string) {
	println('Processing $filename ...')
	lines := os.read_lines(filename) or { panic(err) }
	mut re := regex.new()
	re.compile_opt(r'^(\d+)-(\d+) ([a-z]): (\S+)$') or { panic(err) }
	mut count := 0
	for line in lines {
		re.match_string(line)
		start := re_int_field(re, line, 0)
		end := re_int_field(re, line, 1)
		letter := re_field(re, line, 2)
		passwd := re_field(re, line, 3)
		// println('$line: start=$start, end=$end, letter=$letter, passwd=$passwd')
		if valid(start, end, letter, passwd) {
			count++
		}
	}
	println('$count valid passwords')
}

fn valid(start int, end int, letter string, passwd string) bool {
	first := passwd[start - 1..start] == letter
	second := passwd[end - 1..end] == letter
	return first != second
}

fn re_int_field(re regex.RE, line string, n int) int {
	f := re_field(re, line, n)
	return strconv.atoi(f)
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
