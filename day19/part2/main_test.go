package main

import (
	"log"
	"testing"
)

// 153 too low

func TestMatch(t *testing.T) {
	rules := map[int]*node{
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

	tests := []struct {
		msg       string
		rule      int
		lastMatch int
		length    int
		want      bool
	}{
		// {
		// 	msg:  "aaaabbaaaabbaaa",
		// 	rule: 0,
		// 	want: false,
		// },
		// {
		// 	msg:  "baaaabbaaa",
		// 	rule: 11,
		// 	want: false,
		// },

		// {
		// 	msg:    "aaaaabbaabaaaaababaa",
		// 	rule:   0,
		// 	length: 20,
		// 	want:   true,
		// },

		// {
		// 	msg:    "bbabbbbaabaabba",
		// 	rule:   0,
		// 	length: 15,
		// 	want:   true,
		// },

		{
			msg:  "babbbbaabbbbbabbbbbbaabaaabaaa",
			rule: 0,
			want: true,
		},

		// {
		// 	msg:       "aabaaabaaa",
		// 	rule:      11,
		// 	lastMatch: 42,
		// 	want:      true,
		// },

		// {
		// 	msg:       "aabaaabaaa",
		// 	rule:      11,
		// 	lastMatch: 8,
		// 	length:    10,
		// 	want:      true,
		// },

		//{
		//	msg: "aaabbbbbbaaaabaababaabababbabaaabbababababaaa",
		//	rule: 0,
		//	length:,
		//	want:true,
		//},
		//{
		//	msg: "bbbbbbbaaaabbbbaaabbabaaa",
		//	rule: 0,
		//	length:,
		//	want:true,
		//},
		//{
		//	msg: "bbbababbbbaaaaaaaabbababaaababaabab",
		//	rule: 0,
		//	length:,
		//	want:true,
		//},
		//{
		//	msg: "ababaaaaaabaaab",
		//	rule: 0,
		//	length:,
		//	want:true,
		//},
		//{
		//	msg: "ababaaaaabbbaba",
		//	rule: 0,
		//	length:,
		//	want:true,
		//},
		//{
		//	msg: "baabbaaaabbaaaababbaababb",
		//	rule: 0,
		//	length:,
		//	want:true,
		//},
		//{
		//	msg: "abbbbabbbbaaaababbbbbbaaaababb",
		//	rule: 0,
		//	length:,
		//	want:true,
		//},
		//{
		//	msg: "aaaaabbaabaaaaababaa",
		//	rule: 0,
		//	length:,
		//	want:true,
		//},
		//{
		//	msg: "aaaabbaabbaaaaaaabbbabbbaaabbaabaaa",
		//	rule: 0,
		//	length:,
		//	want:true,
		//},
		//{
		//	msg: "aabbbbbaabbbaaaaaabbbbbababaaaaabbaaabba",
		//	rule: 0,
		//	length:,
		//	want:true,
		//},

		// {
		// 	msg:    "aaaaabbaabaaaaababaa",
		// 	rule:   8,
		// 	length: 20,
		// 	want:   true,
		// },
		// {
		// 	msg:    "aabbaabaaaaababaa",
		// 	rule:   15,
		// 	length: 2,
		// 	want:   true,
		// },
		// {
		// 	msg:       "babaa",
		// 	rule:      11,
		// 	lastMatch: 42,
		// 	length:    5,
		// 	want:      true,
		// },
	}

	*verbose = true
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			n, lm, got := match(tt.msg, tt.rule, tt.lastMatch, rules)
			log.Printf("test result: (n=%v, lm=%v, got=%v)", n, lm, got)

			if got != tt.want {
				t.Errorf("match(%v[%v],%v) = %v, want %v", tt.msg, len(tt.msg), tt.rule, got, tt.want)
			}

			want := tt.length
			if want == 0 {
				want = len(tt.msg)
			}
			if n != want {
				t.Errorf("match(%v[%v],%v) length = %v, want %v", tt.msg, len(tt.msg), tt.rule, n, want)
			}
		})
	}
}
