// -*- compile-command: "v run main.v ../example1.txt ../input.txt"; -*-
import os
import regex
import strconv

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
			if valid_pair(key_val[0], key_val[1]) {
                // println('valid_pair(${key_val[0]}, ${key_val[1]}) = true')
				p.required++
			}
		}
        // println('$passport: ${p.required}, ${required.len}')
		if p.required == required.len {
			result << p
		}
	}
	return result
}

fn valid_pair(key string, val string) bool {
    if !required[key] {
        return false
    }
    match key {
        'byr' {
            if valid_num(val, 1920, 2002) {
                return true
            }
        }
        'iyr' {
            if valid_num(val, 2010, 2020) {
                return true
            }
        }
        'eyr' {
            if valid_num(val, 2020, 2030) {
                return true
            }
        }
        'hgt' {
            if val.ends_with('cm') && valid_num(val[..val.len-2], 150, 193) {
                return true
            } else if val.ends_with('in') && valid_num(val[..val.len-2], 59, 76) {
                return true
            }
        }
        'hcl' {
            if valid_rgb(val) {
                return true
            }
        }
        'ecl' {
            match val {
                'amb', 'blu', 'brn', 'gry', 'grn', 'hzl', 'oth' {
                    return true
                }
                else { }
            }
        }
        'pid' {
            if valid_pid(val) {
                return true
            }
        }
        else { }
    }
    return false
}

fn valid_num(s string, min int, max int) bool {
	v := strconv.atoi(s) or { panic(err) }
    // println('valid_num($s, $min, $max) = ${v >= min && v <= max}')
    return v >= min && v <= max
}

fn valid_rgb(s string) bool {
    mut re := regex.new()
	re.compile_opt(r'^#[0-9a-f]{6}$') or { panic(err) }
    start, _ := re.match_string(s)
    // println('valid_rgb($s) = ${start >= 0}')
    return start >= 0
}

fn valid_pid(s string) bool {
    mut re := regex.new()
    re.compile_opt(r'^\d{9}$') or { panic(err) }
    start, _ := re.match_string(s)
    // println('valid_pid($s) = ${start >= 0}')
    return start >= 0
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
