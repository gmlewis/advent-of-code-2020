// -*- compile-command: "go run main.go ../example1.txt ../input.txt"; -*-

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "Verbose log messages")
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		process(arg)
	}

	log.Printf("Done.")
}

func process(filename string) {
	log.Printf("Processing %v ...", filename)
	buf, err := ioutil.ReadFile(filename)
	check("ReadFile: %v", err)

	tilesStr := strings.Split(string(buf), "\n\n")
	tiles := map[int]*tileT{}
	signatures := map[string]map[int]string{}
	for _, tileStr := range tilesStr {
		tileStr = strings.TrimSpace(tileStr)
		if tileStr == "" {
			continue
		}
		tile := parseTile(tileStr, signatures)
		tiles[tile.id] = tile
	}

	deleteDupedSigs(signatures)

	for k, v := range signatures {
		if len(v) < 2 {
			logf("signature[%v] = %v", k, v)
			continue
		}
		if len(v) > 2 {
			log.Fatalf("signature[%v] = %v", k, v)
		}

		log.Printf("signature[%v] = %v", k, v)

		keys := []int{}
		vals := []string{}
		for k2, v2 := range v {
			keys = append(keys, k2)
			vals = append(vals, v2)
		}

		tiles[keys[0]].connections[k] = connector{fromSide: vals[0], toID: keys[1], toSide: vals[1]}
		tiles[keys[1]].connections[k] = connector{fromSide: vals[1], toID: keys[0], toSide: vals[0]}
	}

	for k, v := range tiles {
		if len(v.connections) == 2 {
			log.Printf("corners[%v] = %v", k, v.connections)
		} else if len(v.connections) == 3 {
			log.Printf("edge[%v] = %v", k, v.connections)
		} else if len(v.connections) == 4 {
			log.Printf("inner[%v] = %v", k, v.connections)
		} else {
			log.Fatalf("unknown[%v] = %v", k, v.connections)
		}
	}

	puz := &puzT{tiles: map[key]*tileT{}}
	for k, v := range tiles {
		if len(v.connections) == 2 && isUpperLeft(v) {
			log.Printf("UPPER LEFT: corners[%v] = %v", k, v.connections)
			puz.tiles[key{x: 0, y: 0}] = v
			break
		}
	}
}

type connector struct {
	fromSide string
	toID     int
	toSide   string
}

func isUpperLeft(t *tileT) bool {
	for k := range t.connections {
		if k == t.nSig || k == t.nRevSig || k == t.wSig || k == t.wRevSig {
			return false
		}
	}
	return true
}

func renderPuzzle(corner int, tiles map[int]*tileT, signatures map[string]map[int]string) *puzT {
	t := tiles[corner]
	puz := &puzT{tiles: map[key]*tileT{}}
	log.Printf("puz=%v", puz)

	var north string
	var south string
	var east string
	var west string
	for k := range t.connections {
		if k == t.nSig || k == t.nRevSig {
			north = k
		}
		if k == t.sSig || k == t.sRevSig {
			south = k
		}
		if k == t.eSig || k == t.eRevSig {
			east = k
		}
		if k == t.wSig || k == t.wRevSig {
			west = k
		}
	}
	log.Printf("north=%v, south=%v, east=%v, west=%v", north, south, east, west)

	dx, dy := 1, 1
	switch {
	case north == "" && west == "":
		log.Printf("UPPER RIGHT")
		t.x, t.y = -1, 0
		dx = -1
	case north == "" && east == "":
		log.Printf("UPPER LEFT")
		t.x, t.y = 0, 0
	case south == "" && west == "":
		log.Printf("LOWER RIGHT")
		t.x, t.y = -1, -1
		dx, dy = -1, -1
	case south == "" && east == "":
		log.Printf("LOWER LEFT")
		t.x, t.y = 0, -1
		dy = -1
	}

	// first corner tile has no rotation or flipping
	puz.tiles[key{x: t.x, y: t.y}] = t
	log.Printf("starting corner: (%v,%v) (%v,%v) %v", t.x, t.y, dx, dy, tiles[corner])

	// seen := map[int]bool{}
	// for len(puz.tiles) < len(tiles) {
	// 	for _, tile := range puz.tiles {
	// 		if seen[tile.id] {
	// 			continue
	// 		}
	// 		seen[tile.id] = true
	// 		xn, yn := mates(tile, dx, dy, signatures, tiles)
	// 		// xn := xNeighbor(tile, dx, signatures, tiles)
	// 		puz.tiles[key{x: xn.x, y: xn.y}] = xn
	// 		// yn := yNeighbor(tile, dy, signatures, tiles)
	// 		puz.tiles[key{x: yn.x, y: yn.y}] = yn
	// 	}
	// }

	return puz
}

func mates(t *tileT, dx, dy int, signatures map[string]map[int]string, tiles map[int]*tileT) (xMate, yMate *tileT) {
	// No orientation changes.
	for _, v := range t.connections {
		switch v.fromSide {
		case "wSig", "wRevSig", "eSig", "eRevSig":
			xMate = tiles[v.toID]
			xMate.x, xMate.y = t.x+dx, t.y
			xMate.setXOrientation(t.orientation, v.fromSide, v.toSide)
		case "nSig", "nRevSig", "sSig", "sRevSig":
			yMate = tiles[v.toID]
			yMate.x, yMate.y = t.x, t.y+dy
			yMate.setYOrientation(t.orientation, v.fromSide, v.toSide)
		default:
			log.Fatalf("unhandled v.fromSize=%v", v.fromSide)
		}
	}
	return xMate, yMate
}

func (t *tileT) setXOrientation(fromOrientation, fromSide, toSide string) {
	key := fmt.Sprintf("%v,%v,%v", fromOrientation, fromSide, toSide)
	switch key {
	case ",wSig,eSig", ",wRevSig,eRevSig":
		t.orientation = ""
	case ",wSig,wRevSig", ",wRevSig,wSig":
		t.orientation = "Rot180"
	case ",eSig,wRevSig", ",eRevSig,wSig", ",wSig,eRevSig", ",wRevSig,eSig":
		t.orientation = "MirrorY"
	case ",eRevSig,wRevSig", ",eSig,wSig":
		t.orientation = ""
	case ",wRevSig,sRevSig", ",wSig,sSig":
		t.orientation = "RotRight"
	default:
		log.Fatalf("setXOrientation: unhandled: case %q:", key)
	}
}

func (t *tileT) setYOrientation(fromOrientation, fromSide, toSide string) {
	key := fmt.Sprintf("%v,%v,%v", fromOrientation, fromSide, toSide)
	switch key {
	case ",sSig,nSig", ",sRevSig,nRevSig", ",nRevSig,sRevSig", ",nSig,sSig":
		t.orientation = ""
	case ",sSig,wRevSig", ",sRevSig,wSig":
		t.orientation = "RotLeft+MirrorX"
	case ",nRevSig,eRevSig", ",nSig,eSig":
		t.orientation = "RotLeft+MirrorX"
	case ",wSig,wRevSig", ",wRevSig,wSig":
		t.orientation = "Rot180"
	default:
		log.Fatalf("setYOrientation: unhandled: case %q:", key)
	}
}

func xNeighbor(t *tileT, dx int, signatures map[string]map[int]string, tiles map[int]*tileT) *tileT {
	// No orientation changes
	for c := range t.connections {
		log.Printf("xNeighbor(%v): checking connection %v ...", t.id, c)
		if c == t.wSig || c == t.wRevSig {
			// west connection
			mate := findMate(t, c, signatures[c], tiles)
			if t.wSig == mate.eSig || t.wRevSig == mate.eRevSig {
				// No orientation changes
				mate.x, mate.y = t.x+dx, t.y
				return mate
			}
			log.Fatalf("unhandled west connection mate: c=%v, %v", c, mate)
		}
		if c == t.eSig || c == t.eRevSig {
			// east connection
			mate := findMate(t, c, signatures[c], tiles)
			if t.eSig == mate.wSig || t.eRevSig == mate.wRevSig {
				// No orientation changes
				mate.x, mate.y = t.x+dx, t.y
				return mate
			}
			log.Fatalf("unhandled east connection mate: c=%v, %v", c, mate)
		}
	}
	log.Fatalf("unhandled case: dx=%v, t=%v", dx, t)
	return nil
}

func yNeighbor(t *tileT, dy int, signatures map[string]map[int]string, tiles map[int]*tileT) *tileT {
	// No orientation changes
	for c := range t.connections {
		log.Printf("yNeighbor: checking connection %v ...", c)
		if c == t.nSig || c == t.nRevSig {
			// north connection
			mate := findMate(t, c, signatures[c], tiles)
			if t.nSig == mate.sSig || t.nRevSig == mate.sRevSig {
				// No orientation changes
				mate.x, mate.y = t.x, t.y+dy
				return mate
			}
			log.Fatalf("unhandled north connection mate: c=%v, %v", c, mate)
		}
		if c == t.sSig || c == t.sRevSig {
			// south connection
			mate := findMate(t, c, signatures[c], tiles)
			if t.sSig == mate.nSig || t.sRevSig == mate.nRevSig {
				// No orientation changes
				mate.x, mate.y = t.x, t.y+dy
				return mate
			}
			log.Fatalf("unhandled south connection mate: c=%v, %v", c, mate)
		}
	}
	log.Fatalf("dy=%v, t=%v", dy, t)
	return nil
}

func findMate(t *tileT, c string, pairs map[int]string, tiles map[int]*tileT) *tileT {
	for k := range pairs {
		if k != t.id {
			return tiles[k]
		}
	}
	log.Fatalf("findMate failed")
	return nil
}

func (p *puzT) searchPuzzle() (monsters, roughWater int) {
	return monsters, roughWater
}

type puzT struct {
	tiles map[key]*tileT
}

type tileT struct {
	id     int
	dots   map[key]bool
	width  int
	height int

	nSig    string
	sSig    string
	eSig    string
	wSig    string
	nRevSig string
	sRevSig string
	eRevSig string
	wRevSig string

	connections map[string]connector

	x           int
	y           int
	orientation string
}

func parseTile(s string, signatures map[string]map[int]string) *tileT {
	lines := strings.Split(s, "\n")
	idStr := lines[0][5 : len(lines[0])-1]
	id, err := strconv.Atoi(idStr)
	check("id: %v", err)

	t := &tileT{id: id, dots: map[key]bool{}, connections: map[string]connector{}}
	var eastSig string
	var westSig string
	var southSig string
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if t.height == 0 {
			t.addSignature("nSig", line, signatures)
			t.nSig = line
			t.nRevSig = reverse(line)
			t.addSignature("nRevSig", t.nRevSig, signatures)
		}
		westSig += line[0:1]
		eastSig += line[len(line)-1 : len(line)]
		southSig = line

		t.width = len(line)
		for x := 0; x < len(line); x++ {
			if line[x:x+1] == "#" {
				t.dots[key{x: x, y: t.height}] = true
			}
		}
		t.height++
	}

	t.addSignature("eSig", eastSig, signatures)
	t.eSig = eastSig
	t.eRevSig = reverse(eastSig)
	t.addSignature("eRevSig", t.eRevSig, signatures)

	t.addSignature("wSig", westSig, signatures)
	t.wSig = westSig
	t.wRevSig = reverse(westSig)
	t.addSignature("wRevSig", t.wRevSig, signatures)

	t.addSignature("sSig", southSig, signatures)
	t.sSig = southSig
	t.sRevSig = reverse(southSig)
	t.addSignature("sRevSig", t.sRevSig, signatures)

	log.Printf("parsed:\n%v", t)
	return t
}

func (t *tileT) addSignature(name, sig string, signatures map[string]map[int]string) {
	if _, ok := signatures[sig]; !ok {
		signatures[sig] = map[int]string{}
	}
	signatures[sig][t.id] = name
}

func (t *tileT) String() string {
	lines := []string{fmt.Sprintf("Tile %v:", t.id)}
	lines = append(lines, fmt.Sprintf("nSig: %v nRevSig: %v", t.nSig, t.nRevSig))
	lines = append(lines, fmt.Sprintf("sSig: %v sRevSig: %v", t.sSig, t.sRevSig))
	lines = append(lines, fmt.Sprintf("eSig: %v eRevSig: %v", t.eSig, t.eRevSig))
	lines = append(lines, fmt.Sprintf("wSig: %v wRevSig: %v", t.wSig, t.wRevSig))
	lines = append(lines, fmt.Sprintf("Position: (%v,%v)", t.x, t.y))
	lines = append(lines, fmt.Sprintf("Connections: %v", t.connections))
	for y := 0; y < t.height; y++ {
		var line string
		for x := 0; x < t.width; x++ {
			if t.dots[key{x: x, y: y}] {
				line += "#"
			} else {
				line += "."
			}
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

type key struct {
	x int
	y int
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func deleteDupedSigs(signatures map[string]map[int]string) {
	toDelete := map[string]bool{}
	for k, v := range signatures {
		rev := reverse(k)
		if rev == k || len(v) != len(signatures[rev]) || toDelete[k] || toDelete[rev] {
			continue
		}

		keys := map[int]bool{}
		for k2 := range v {
			keys[k2] = true
		}

		keysRev := map[int]bool{}
		for k2 := range signatures[rev] {
			keysRev[k2] = true
		}

		if reflect.DeepEqual(keys, keysRev) {
			toDelete[rev] = true
		}
	}
	for k := range toDelete {
		logf("DELETING duped signature[%v] = %v", k, signatures[k])
		delete(signatures, k)
	}
}

func check(fmtStr string, args ...interface{}) {
	if err := args[len(args)-1]; err != nil {
		log.Fatalf(fmtStr, args...)
	}
}

func logf(fmt string, args ...interface{}) {
	if *verbose {
		log.Printf(fmt, args...)
	}
}
