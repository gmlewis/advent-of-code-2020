package main

import (
	"log"
	"testing"
)

// 153 too low

func TestMatch(t *testing.T) {
	tests := []struct {
		msg    string
		seq    []int
		want   int
		wantOK bool
	}{
		{
			msg:    "",
			want:   0,
			wantOK: true,
		},
		{
			msg:    "",
			seq:    []int{},
			want:   0,
			wantOK: true,
		},
		{
			msg:    "",
			seq:    []int{1},
			want:   0,
			wantOK: false,
		},
		{
			msg:    "a",
			seq:    []int{},
			want:   0,
			wantOK: false,
		},
		{
			msg:    "b",
			seq:    []int{},
			want:   0,
			wantOK: false,
		},
		{
			msg:    "a",
			seq:    []int{1},
			want:   1,
			wantOK: true,
		},
		{
			msg:    "b",
			seq:    []int{1},
			want:   0,
			wantOK: false,
		},
		{
			msg:    "b",
			seq:    []int{14},
			want:   1,
			wantOK: true,
		},
		{
			msg:    "a",
			seq:    []int{15},
			want:   1,
			wantOK: true,
		},
		{
			msg:    "b",
			seq:    []int{15},
			want:   1,
			wantOK: true,
		},
		{
			msg:    "aa",
			seq:    []int{18},
			want:   2,
			wantOK: true,
		},
		{
			msg:    "ab",
			seq:    []int{18},
			want:   2,
			wantOK: true,
		},
		{
			msg:    "ba",
			seq:    []int{18},
			want:   2,
			wantOK: true,
		},
		{
			msg:    "bb",
			seq:    []int{18},
			want:   2,
			wantOK: true,
		},
		{
			msg:    "aaaaabbaabaaaaababaa",
			seq:    []int{0},
			want:   len("aaaaabbaabaaaaababaa"),
			wantOK: true,
		},
		// {
		// 	msg:  "aaaabbaaaabbaaa",
		// 	seq: []int{0},
		// want: false,
		// },
		// {
		// 	msg:  "aaaabbaabbaaaaaaabbbabbbaaabbaabaaa",
		// 	seq: []int{0},
		// want: true,
		// },
		// {
		// 	msg:  "aaabbbbbbaaaabaababaabababbabaaabbababababaaa",
		// 	seq: []int{0},
		// want: true,
		// },
		// {
		// 	msg:  "aabbbbbaabbbaaaaaabbbbbababaaaaabbaaabba",
		// 	seq: []int{0},
		// want: true,
		// },
		// {
		// 	msg:  "ababaaaaaabaaab",
		// 	seq: []int{0},
		// want: true,
		// },
		// {
		// 	msg:  "ababaaaaabbbaba",
		// 	seq: []int{0},
		// want: true,
		// },
		// {
		// 	msg:  "abbbbabbbbaaaababbbbbbaaaababb",
		// 	seq: []int{0},
		// want: true,
		// },
		// {
		// 	msg:  "abbbbbabbbaaaababbaabbbbabababbbabbbbbbabaaaa",
		// 	seq: []int{0},
		// want: false,
		// },
		// {
		// 	msg:  "baabbaaaabbaaaababbaababb",
		// 	seq: []int{0},
		// want: true,
		// },
		// {
		// 	msg:  "babaaabbbaaabaababbaabababaaab",
		// 	seq: []int{0},
		// want: false,
		// },
		// {
		// 	msg:  "babbbbaabbbbbabbbbbbaabaaabaaa",
		// 	seq: []int{0},
		// want: true,
		// },
		// {
		// 	msg:  "bbabbbbaabaabba",
		// 	seq: []int{0},
		// want: true,
		// },
		// {
		// 	msg:  "bbbababbbbaaaaaaaabbababaaababaabab",
		// 	seq: []int{0},
		// want: true,
		// },
		// {
		// 	msg:  "bbbbbbbaaaabbbbaaabbabaaa",
		// 	seq: []int{0},
		// want: true,
		// },
	}

	*verbose = true
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			log.Printf("TEST: match(%q[%v],%v)", tt.msg, len(tt.msg), tt.seq)
			n, ok := match(tt.msg, tt.seq, rules)

			if n != tt.want {
				t.Errorf("match(%v[%v],%v) = %v, want %v", tt.msg, len(tt.msg), tt.seq, n, tt.want)
			}
			if ok != tt.wantOK {
				t.Errorf("match(%v[%v],%v) = %v, want %v", tt.msg, len(tt.msg), tt.seq, ok, tt.wantOK)
			}
		})
	}
}

var rules = map[int]*node{
	0:  &node{matches: map[string]bool{}, s: "", seq1: []int{8, 11}, seq2: []int(nil)},
	1:  &node{matches: map[string]bool{"a": true}, s: "a", seq1: []int(nil), seq2: []int(nil)},
	2:  &node{matches: map[string]bool{}, s: "", seq1: []int{1, 24}, seq2: []int{14, 4}},
	3:  &node{matches: map[string]bool{}, s: "", seq1: []int{5, 14}, seq2: []int{16, 1}},
	4:  &node{matches: map[string]bool{}, s: "", seq1: []int{1, 1}, seq2: []int(nil)},
	5:  &node{matches: map[string]bool{}, s: "", seq1: []int{1, 14}, seq2: []int{15, 1}},
	6:  &node{matches: map[string]bool{}, s: "", seq1: []int{14, 14}, seq2: []int{1, 14}},
	7:  &node{matches: map[string]bool{}, s: "", seq1: []int{14, 5}, seq2: []int{1, 21}},
	8:  &node{matches: map[string]bool{}, s: "", seq1: []int{42}, seq2: []int{42, 8}},
	9:  &node{matches: map[string]bool{}, s: "", seq1: []int{14, 27}, seq2: []int{1, 26}},
	10: &node{matches: map[string]bool{}, s: "", seq1: []int{23, 14}, seq2: []int{28, 1}},
	11: &node{matches: map[string]bool{}, s: "", seq1: []int{42, 31}, seq2: []int{42, 11, 31}},
	12: &node{matches: map[string]bool{}, s: "", seq1: []int{24, 14}, seq2: []int{19, 1}},
	13: &node{matches: map[string]bool{}, s: "", seq1: []int{14, 3}, seq2: []int{1, 12}},
	14: &node{matches: map[string]bool{"b": true}, s: "b", seq1: []int(nil), seq2: []int(nil)},
	15: &node{matches: map[string]bool{}, s: "", seq1: []int{1}, seq2: []int{14}},
	16: &node{matches: map[string]bool{}, s: "", seq1: []int{15, 1}, seq2: []int{14, 14}},
	17: &node{matches: map[string]bool{}, s: "", seq1: []int{14, 2}, seq2: []int{1, 7}},
	18: &node{matches: map[string]bool{}, s: "", seq1: []int{15, 15}, seq2: []int(nil)},
	19: &node{matches: map[string]bool{}, s: "", seq1: []int{14, 1}, seq2: []int{14, 14}},
	20: &node{matches: map[string]bool{}, s: "", seq1: []int{14, 14}, seq2: []int{1, 15}},
	21: &node{matches: map[string]bool{}, s: "", seq1: []int{14, 1}, seq2: []int{1, 14}},
	22: &node{matches: map[string]bool{}, s: "", seq1: []int{14, 14}, seq2: []int(nil)},
	23: &node{matches: map[string]bool{}, s: "", seq1: []int{25, 1}, seq2: []int{22, 14}},
	24: &node{matches: map[string]bool{}, s: "", seq1: []int{14, 1}, seq2: []int(nil)},
	25: &node{matches: map[string]bool{}, s: "", seq1: []int{1, 1}, seq2: []int{1, 14}},
	26: &node{matches: map[string]bool{}, s: "", seq1: []int{14, 22}, seq2: []int{1, 20}},
	27: &node{matches: map[string]bool{}, s: "", seq1: []int{1, 6}, seq2: []int{14, 18}},
	28: &node{matches: map[string]bool{}, s: "", seq1: []int{16, 1}, seq2: []int(nil)},
	31: &node{matches: map[string]bool{}, s: "", seq1: []int{14, 17}, seq2: []int{1, 13}},
	42: &node{matches: map[string]bool{}, s: "", seq1: []int{9, 14}, seq2: []int{10, 1}},
}
