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

		logf("signature[%v] = %v", k, v)

		for k2 := range v {
			tiles[k2].connections[k] = true
		}
	}

	corners := map[int]bool{}
	prod := int64(1)
	for k, v := range tiles {
		if len(v.connections) == 2 {
			logf("corners[%v] = %v", k, v.connections)
			corners[k] = true
			prod *= int64(k)
		}
	}
	log.Printf("Solution: %v", prod)
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

	connections map[string]bool
}

func parseTile(s string, signatures map[string]map[int]string) *tileT {
	lines := strings.Split(s, "\n")
	idStr := lines[0][5 : len(lines[0])-1]
	id, err := strconv.Atoi(idStr)
	check("id: %v", err)

	t := &tileT{id: id, dots: map[key]bool{}, connections: map[string]bool{}}
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
		eastSig += line[0:1]
		westSig += line[len(line)-1 : len(line)]
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

	logf("parsed:\n%v", t)
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
