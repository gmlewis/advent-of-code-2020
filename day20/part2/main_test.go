package main

import (
	"reflect"
	"testing"
)

func TestMirrorX(t *testing.T) {
	sigs1 := map[string]map[int]string{}
	tile := parseTile(`Tile 0:
#.#.#####.
.#..######
..#.......
######....
####.#..#.
.#...#.##.
#.#####.##
..#.###...
..#.......
..#.###...`, sigs1)

	sigs2 := map[string]map[int]string{}
	want := parseTile(`Tile 0:
.#####.#.#
######..#.
.......#..
....######
.#..#.####
.##.#...#.
##.#####.#
...###.#..
.......#..
...###.#..`, sigs2)
	want.nSig = ""
	want.nRevSig = ""
	want.sSig = ""
	want.sRevSig = ""
	want.wSig = ""
	want.wRevSig = ""
	want.eSig = ""
	want.eRevSig = ""

	got := mirrorX(tile)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mirrorX =\n%v\nwant:\n%v", got, want)
	}
}

func TestMirrorY(t *testing.T) {
	sigs1 := map[string]map[int]string{}
	tile := parseTile(`Tile 0:
#.#.#####.
.#..######
..#.......
######....
####.#..#.
.#...#.##.
#.#####.##
..#.###...
..#.......
..#.###...`, sigs1)

	sigs2 := map[string]map[int]string{}
	want := parseTile(`Tile 0:
..#.###...
..#.......
..#.###...
#.#####.##
.#...#.##.
####.#..#.
######....
..#.......
.#..######
#.#.#####.`, sigs2)
	want.nSig = ""
	want.nRevSig = ""
	want.sSig = ""
	want.sRevSig = ""
	want.wSig = ""
	want.wRevSig = ""
	want.eSig = ""
	want.eRevSig = ""

	got := mirrorY(tile)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mirrorY =\n%v\nwant:\n%v", got, want)
	}
}

/*
(defun my-reverse-region (beg end)
 "Reverse characters between BEG and END."
 (interactive "r")
 (let ((region (buffer-substring beg end)))
   (delete-region beg end)
   (insert (nreverse region))))
*/
