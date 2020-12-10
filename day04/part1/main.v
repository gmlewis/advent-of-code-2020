// -*- compile-command: "v run main.v ../example1.txt ../input.txt"; -*-
import os

fn main() {
	for arg in os.args[1..] {
		process(arg)
	}
	println('Done.')
}

struct Passport {
mut:
	required int
	pairs    map[string]string
}

fn process(filename string) {
	println('Processing $filename ...')
	passports := read_passports(filename)
	println('$passports.len valid passports')
}

fn read_passports(filename string) []&Passport {
	mut result := []&Passport{}
	buf := os.read_file(filename) or { panic(err) }
	passports := buf.split('\n\n')
	for passport in passports {
		parts := passport.replace('\n', ' ').trim_space().split(' ')
		mut p := &Passport{}
		for part in parts {
			key_val := part.split(':')
			p.pairs[key_val[0]] = key_val[1]
			if required[key_val[0]] {
				p.required++
			}
		}
		if p.required == required.len {
			result << p
		}
	}
	return result
}

const (
    required = {
	'byr': true, // (Birth Year)
	'iyr': true, // (Issue Year)
	'eyr': true, // (Expiration Year)
	'hgt': true, // (Height)
	'hcl': true, // (Hair Color)
	'ecl': true, // (Eye Color)
	'pid': true, // (Passport ID)
	// 'cid' // (Country ID)
    }
)
