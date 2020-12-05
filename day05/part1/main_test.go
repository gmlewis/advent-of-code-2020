package main

import (
	"fmt"
	"testing"
)

func TestSpaceID(t *testing.T) {
	tests := []struct {
		space string
		want  int
	}{
		{space: "FBFBBFFRLR", want: 357},
		{space: "BFFFBBFRRR", want: 567},
		{space: "FFFBBBFRRR", want: 119},
		{space: "BBFFBBFRLL", want: 820},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%v", i), func(t *testing.T) {
			got := spaceID(tt.space)

			if got != tt.want {
				t.Errorf("spaceID(%q) = %v, want %v", tt.space, got, tt.want)
			}
		})
	}
}
