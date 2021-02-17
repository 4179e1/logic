package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

const nonPrefix = "non-"

var depthFlag int
var cache = map[string]string{}

func init() {

	flag.IntVar(&depthFlag, "d", 16, "Depth")
}

func tabPrint(depthStr, str string) int {
	depth := len(strings.Split(depthStr, ".")) - 1

	for i := 0; i < depth; i++ {
		fmt.Print("\t")
	}
	fmt.Printf("%s. %s\n", depthStr, str)

	return depth
}

func non(x string) string {
	if strings.HasPrefix(x, nonPrefix) {
		return strings.TrimPrefix(x, nonPrefix)
	}
	return nonPrefix + x
}

func astr(s, p string) string {
	return fmt.Sprintf("All %s is %s", s, p)
}

func A(s, p string, depthStr string) {
	str := astr(s, p)
	idx, found := cache[str]
	if found {
		tabPrint(depthStr, fmt.Sprintf("%s [Dup %s]", str, idx))
		return
	}
	cache[str] = depthStr

	depth := tabPrint(depthStr, str)
	if depth >= depthFlag {
		return
	}

	// Converse
	I(p, s, depthStr+".1")

	// Obverse
	E(s, non(p), depthStr+".2")

	// Conrapositive
	A(non(p), non(s), depthStr+".3")

}

func estr(s, p string) string {
	return fmt.Sprintf("No %s is %s", s, p)
}

func E(s, p string, depthStr string) {
	str := estr(s, p)
	idx, found := cache[str]
	if found {
		tabPrint(depthStr, fmt.Sprintf("%s [Dup %s]", str, idx))
		return
	}
	cache[str] = depthStr

	depth := tabPrint(depthStr, str)
	if depth >= depthFlag {
		return
	}

	// Converse
	E(p, s, depthStr+".1")

	// Obverse
	A(s, non(p), depthStr+".2")

	// Conrapositive
	O(non(p), non(s), depthStr+".3")

}

func istr(s, p string) string {
	return fmt.Sprintf("Some %s is %s", s, p)
}

func I(s, p string, depthStr string) {
	str := istr(s, p)

	idx, found := cache[str]
	if found {
		tabPrint(depthStr, fmt.Sprintf("%s [Dup %s]", str, idx))
		return
	}
	cache[str] = depthStr

	depth := tabPrint(depthStr, str)
	if depth >= depthFlag {
		return
	}

	// Converse
	I(p, s, depthStr+".1")

	// Obverse
	O(s, non(p), depthStr+".2")

	// Conrapositive
	// TODO
	tabPrint(depthStr+".3", "[Contraposition not valid]")

}

func ostr(s, p string) string {
	return fmt.Sprintf("Some %s is not %s", s, p)
}
func O(s, p string, depthStr string) {
	str := ostr(s, p)

	idx, found := cache[str]
	if found {
		tabPrint(depthStr, fmt.Sprintf("%s [Dup %s]", str, idx))
		return
	}
	cache[str] = depthStr

	depth := tabPrint(depthStr, str)
	if depth >= depthFlag {
		return
	}

	// Converse
	// TODO
	tabPrint(depthStr+".1", "[Conversion not valid]")
	//A(s, p, depthStr+".1")

	// Obverse
	I(s, non(p), depthStr+".2")

	// Conrapositive
	O(non(p), non(s), depthStr+".3")

}

func opposition(prop, s, p string) {
	switch prop {
	case "A":
		A(s, p, "0")
	case "E":
		E(s, p, "0")
	case "I":
		I(s, p, "0")
	case "O":
		O(s, p, "0")
	default:
		fmt.Fprintf(os.Stderr, "Unknown Proposition %s\n", prop)
		os.Exit(2)
	}

	rch := map[string]string{}
	idxs := []string{}
	for k, v := range cache {
		idxs = append(idxs, v)
		rch[v] = k
		//fmt.Printf("%-16s %s\n", v, k)
	}

	fmt.Printf("\n=== %d Valid Propositions ===\n", len(cache))
	sort.Strings(idxs)
	for _, idx := range idxs {
		fmt.Printf("%-16s %s\n", idx, rch[idx])
	}
}

func main() {

	flag.Parse()
	args := flag.Args()
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <Proposition> <Subject> <Preditate>\n", os.Args[0])
		os.Exit(1)
	}

	opposition(args[0], args[1], args[2])

}
