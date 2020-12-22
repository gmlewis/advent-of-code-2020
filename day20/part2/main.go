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

		tiles[keys[0]].connections[k] = &connector{fromSide: vals[0], toID: keys[1], toSide: vals[1]}
		tiles[keys[1]].connections[k] = &connector{fromSide: vals[1], toID: keys[0], toSide: vals[0]}
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
		if len(v.connections) == 2 && isLowerLeft(v) {
			log.Printf("LOWER LEFT: corners[%v] = %v", k, v.connections)
			puz.tiles[key{x: 0, y: 0}] = v
			break
		}
	}

	puz.findRightAndUp(key{x: 0, y: 0}, tiles)

	image := puz.render()
	log.Printf("image: %v", image)
}

func (p *puzT) render() *tileT {
	image := &tileT{dots: map[key]bool{}}
	for k, t := range p.tiles {
		log.Printf("Rendering tile %v (%v,%v)", k, t.width, t.height)
		for k2 := range t.dots {
			if k2.x == 0 && k2.y == 0 || k2.x == t.width-1 || k2.y == t.height-1 {
				continue // remove border
			}
			x := k.x*(t.width-2) + k2.x - 1
			y := k.y*(t.height-2) + k2.y - 1
			image.dots[key{x: x, y: y}] = true
			if x >= image.width {
				image.width = x + 1
			}
			if y >= image.height {
				image.height = y + 1
			}
		}
	}
	return image
}

func (p *puzT) findRightAndUp(k key, tiles map[int]*tileT) {
	t := p.tiles[k]
	if t == nil {
		return
	}

	upKey := key{x: k.x, y: k.y + 1}
	rightKey := key{x: k.x + 1, y: k.y}
	for _, v := range t.connections {
		switch v.fromSide {
		case "nSig", "nRevSig":
			if _, ok := p.tiles[upKey]; ok {
				continue
			}
			up := tiles[v.toID]
			log.Printf("BEFORE: found up from %v to %v: %v", v.fromSide, v.toSide, up)
			p.tiles[upKey] = rotateAndFlip(v.fromSide, v.toSide, up)
			log.Printf("AFTER: found up from %v to %v: %v", v.fromSide, v.toSide, p.tiles[upKey])
		case "eSig", "eRevSig":
			if _, ok := p.tiles[rightKey]; ok {
				continue
			}
			right := tiles[v.toID]
			log.Printf("BEFORE: found right from %v to %v: %v", v.fromSide, v.toSide, right)
			p.tiles[rightKey] = rotateAndFlip(v.fromSide, v.toSide, right)
			log.Printf("AFTER: found right from %v to %v: [%v]=%v", v.fromSide, v.toSide, rightKey, p.tiles[rightKey])
		default:
			// log.Printf("key=%v, t.id=%v, unhandled fromSig: %v", k, t.id, v)
		}
	}

	p.findRightAndUp(upKey, tiles)
	p.findRightAndUp(rightKey, tiles)
}

func rotateAndFlip(fromSide, toSide string, t *tileT) *tileT {
	fromTo := fmt.Sprintf("%v,%v", fromSide, toSide)
	switch fromTo {
	case "eSig,eRevSig", "eRevSig,eSig", "wSig,wRevSig", "wRevSig,wSig", "nRevSig,nSig", "nSig,nRevSig", "sRevSig,sSig", "sSig,sRevSig":
		return rot180(t)
	case "nSig,wSig", "nRevSig,wRevSig":
		return rotLeft(t)
	case "eSig,wSig", "eRevSig,wRevSig", "nSig,sSig", "nRevSig,sRevSig":
		return t
	case "eRevSig,sRevSig", "eSig,sSig":
		return rotRight(t)
	case "eSig,wRevSig", "eRevSig,wSig":
		return mirrorY(t)
	case "eRevSig,eRevSig", "eSig,eSig", "wRevSig,wRevSig", "wSig,wSig", "nSig,sRevSig", "nRevSig,sSig", "sSig,nRevSig", "sRevSig,nSig":
		return mirrorX(t)
	default:
		log.Fatalf("unhandled fromTo: case %q:", fromTo)
	}
	return t
}

func mirrorX(t *tileT) *tileT {
	cp := t.copyTile()
	xform(cp, t, func(k key) key {
		return key{x: t.width - 1 - k.x, y: k.y}
	})
	log.Printf("BEFORE: mirrorX: %v", cp.connections)
	for _, v := range cp.connections {
		switch v.fromSide {
		case "sSig":
			v.fromSide = "sRevSig"
		case "sRevSig":
			v.fromSide = "sSig"
		case "nSig":
			v.fromSide = "nRevSig"
		case "nRevSig":
			v.fromSide = "nSig"
		case "wSig":
			v.fromSide = "eSig"
		case "wRevSig":
			v.fromSide = "eRevSig"
		case "eSig":
			v.fromSide = "wSig"
		case "eRevSig":
			v.fromSide = "wRevSig"
		default:
			log.Fatalf("mirrorX: unhandled case %q:", v.fromSide)
		}
	}
	log.Printf("AFTER: mirrorX: %v", cp.connections)
	return cp
}

func mirrorY(t *tileT) *tileT {
	cp := t.copyTile()
	xform(cp, t, func(k key) key {
		return key{x: k.x, y: t.height - 1 - k.y}
	})
	log.Printf("BEFORE: mirrorY: %v", cp.connections)
	for _, v := range cp.connections {
		switch v.fromSide {
		case "sSig":
			v.fromSide = "nSig"
		case "sRevSig":
			v.fromSide = "nRevSig"
		case "nSig":
			v.fromSide = "sSig"
		case "nRevSig":
			v.fromSide = "sRevSig"
		case "wSig":
			v.fromSide = "wRevSig"
		case "wRevSig":
			v.fromSide = "wSig"
		case "eSig":
			v.fromSide = "eRevSig"
		case "eRevSig":
			v.fromSide = "eSig"
		default:
			log.Fatalf("mirrorY: unhandled case %q:", v.fromSide)
		}
	}
	log.Printf("AFTER: mirrorY: %v", cp.connections)
	return cp
}

func rot180(t *tileT) *tileT {
	cp := t.copyTile()
	xform(cp, t, func(k key) key {
		return key{x: t.width - 1 - k.x, y: t.height - 1 - k.y}
	})
	log.Printf("BEFORE: rot180: %v", cp.connections)
	for _, v := range cp.connections {
		switch v.fromSide {
		case "sSig":
			v.fromSide = "nRevSig"
		case "sRevSig":
			v.fromSide = "nSig"
		case "nSig":
			v.fromSide = "sRevSig"
		case "nRevSig":
			v.fromSide = "sSig"
		case "wSig":
			v.fromSide = "eRevSig"
		case "wRevSig":
			v.fromSide = "eSig"
		case "eSig":
			v.fromSide = "wRevSig"
		case "eRevSig":
			v.fromSide = "wSig"
		default:
			log.Fatalf("rot180: unhandled case %q:", v.fromSide)
		}
	}
	log.Printf("AFTER: rot180: %v", cp.connections)
	return cp
}

func rotLeft(t *tileT) *tileT {
	cp := t.copyTile()
	xform(cp, t, func(k key) key {
		return key{x: t.height - 1 - k.y, y: k.x}
	})
	log.Printf("BEFORE: rotLeft: %v", cp.connections)
	for _, v := range cp.connections {
		switch v.fromSide {
		case "wSig":
			v.fromSide = "sSig"
		case "wRevSig":
			v.fromSide = "sRevSig"
		case "sSig":
			v.fromSide = "eRevSig"
		case "sRevSig":
			v.fromSide = "eSig"
		case "eSig":
			v.fromSide = "nSig"
		case "eRevSig":
			v.fromSide = "nRevSig"
		case "nSig":
			v.fromSide = "wRevSig"
		case "nRevSig":
			v.fromSide = "wSig"
		default:
			log.Fatalf("rotLeft: unhandled case %q:", v.fromSide)
		}
	}
	log.Printf("AFTER: rotLeft: %v", cp.connections)
	return cp
}

func rotRight(t *tileT) *tileT {
	cp := t.copyTile()
	xform(cp, t, func(k key) key {
		return key{x: t.height - 1 - k.y, y: k.x}
	})
	log.Printf("BEFORE: rotRight: %v", cp.connections)
	for _, v := range cp.connections {
		switch v.fromSide {
		case "wSig":
			v.fromSide = "nRevSig"
		case "wRevSig":
			v.fromSide = "nSig"
		case "sSig":
			v.fromSide = "wSig"
		case "sRevSig":
			v.fromSide = "eRevSig"
		case "eSig":
			v.fromSide = "sRevSig"
		case "eRevSig":
			v.fromSide = "sSig"
		case "nSig":
			v.fromSide = "eSig"
		case "nRevSig":
			v.fromSide = "eRevSig"
		default:
			log.Fatalf("rotRight: unhandled case %q:", v.fromSide)
		}
	}
	log.Printf("AFTER: rotRight: %v", cp.connections)
	return cp
}

func xform(cp, t *tileT, f func(key) key) {
	for k := range t.dots {
		newKey := f(k)
		// log.Printf("xform %v => %v", k, newKey)
		cp.dots[newKey] = true
	}
}

func (t *tileT) copyTile() *tileT {
	cp := &tileT{
		id:          t.id,
		dots:        map[key]bool{},
		width:       t.width,
		height:      t.height,
		connections: map[string]*connector{},
		x:           t.x,
		y:           t.y,
	}
	for k, v := range t.connections {
		cp.connections[k] = &connector{fromSide: v.fromSide, toID: v.toID, toSide: v.toSide}
	}
	return cp
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

func isLowerLeft(t *tileT) bool {
	for k := range t.connections {
		if k == t.sSig || k == t.sRevSig || k == t.wSig || k == t.wRevSig {
			return false
		}
	}
	return true
}

/*
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
*/

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

	connections map[string]*connector

	x int
	y int
	// orientation string
}

func parseTile(s string, signatures map[string]map[int]string) *tileT {
	lines := strings.Split(s, "\n")
	idStr := lines[0][5 : len(lines[0])-1]
	id, err := strconv.Atoi(idStr)
	check("id: %v", err)

	t := &tileT{id: id, dots: map[key]bool{}, connections: map[string]*connector{}}
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
	lines := []string{fmt.Sprintf("Tile %v (%v,%v):", t.id, t.width, t.height)}
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
