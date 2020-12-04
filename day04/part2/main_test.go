package main

import (
	"fmt"
	"testing"
)

func TestValidPair(t *testing.T) {
	tests := []struct {
		key  string
		val  string
		want bool
	}{
		{key: "byr", want: true, val: "2002"},
		{key: "byr", want: false, val: "2003"},
		{key: "hgt", want: true, val: "60in"},
		{key: "hgt", want: true, val: "190cm"},
		{key: "hgt", want: false, val: "190in"},
		{key: "hgt", want: false, val: "190"},
		{key: "hcl", want: true, val: "#123abc"},
		{key: "hcl", want: false, val: "#123abz"},
		{key: "hcl", want: false, val: "123abc"},
		{key: "ecl", want: true, val: "brn"},
		{key: "ecl", want: false, val: "wat"},
		{key: "pid", want: true, val: "000000001"},
		{key: "pid", want: false, val: "0123456789"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := validPair(tt.key, tt.val)

			if got != tt.want {
				t.Errorf("validPair(%q,%q) = %v, want %v", tt.key, tt.val, got, tt.want)
			}
		})
	}
}
