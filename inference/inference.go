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
var invalidCache = map[string]string{}

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

func inCache(depthStr, str string) bool {
	idx, found := cache[str]
	if found {
		tabPrint(depthStr, fmt.Sprintf("%s [Dup %s]", str, idx))
		return true
	}

	cache[str] = depthStr

	return false
}

func inInvalidCache(depthStr, str, hint string) bool {
	idx, found := invalidCache[str]
	if found {
		tabPrint(depthStr, fmt.Sprintf("%s [Dup %s]%s", str, idx, hint))
		return true
	}

	invalidCache[str] = depthStr + hint

	return false
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

	if inCache(depthStr, str) {
		return
	}

	depth := tabPrint(depthStr, str)
	if depth >= depthFlag {
		return
	}

	// Converse, by limitation
	I(p, s, depthStr+".1L")
	// Converse, Invalid
	invalidStr := astr(p, s)
	invalidIdx := depthStr + ".1X"
	hint := " [Conversion not valid]"
	if !inInvalidCache(invalidIdx, invalidStr, hint) {
		tabPrint(invalidIdx, invalidStr+hint)
	}

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

	if inCache(depthStr, str) {
		return
	}
	/*
		idx, found := cache[str]
		if found {
			tabPrint(depthStr, fmt.Sprintf("%s [Dup %s]", str, idx))
			return
		}
		cache[str] = depthStr
	*/

	depth := tabPrint(depthStr, str)
	if depth >= depthFlag {
		return
	}

	// Converse
	E(p, s, depthStr+".1")

	// Obverse
	A(s, non(p), depthStr+".2")

	// Conrapositive, by limitation
	O(non(p), non(s), depthStr+".3L")
	// Conrapositive, not valid
	invalidStr := estr(non(p), non(s))
	invalidIdx := depthStr + ".3X"
	hint := " [Contraposition not valid]"
	if !inInvalidCache(invalidIdx, invalidStr, hint) {
		tabPrint(invalidIdx, invalidStr+hint)
	}

}

func istr(s, p string) string {
	return fmt.Sprintf("Some %s is %s", s, p)
}

func I(s, p string, depthStr string) {
	str := istr(s, p)

	if inCache(depthStr, str) {
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

	// Conrapositive, not valid
	invalidStr := istr(non(p), non(s))
	invalidIdx := depthStr + ".3"
	hint := " [Contraposition not valid]"
	if !inInvalidCache(invalidIdx, invalidStr, hint) {
		tabPrint(invalidIdx, invalidIdx+hint)
	}

}

func ostr(s, p string) string {
	return fmt.Sprintf("Some %s is not %s", s, p)
}
func O(s, p string, depthStr string) {
	str := ostr(s, p)

	if inCache(depthStr, str) {
		return
	}

	depth := tabPrint(depthStr, str)
	if depth >= depthFlag {
		return
	}

	// Converse, not valid
	invalidStr := ostr(p, s)
	invalidIdx := depthStr + ".1"
	hint := " [Conversion not valid]"
	if !inInvalidCache(invalidIdx, invalidStr, hint) {
		tabPrint(invalidIdx, invalidStr+hint)
	}

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
		fmt.Printf("%-20s %s\n", idx, rch[idx])
	}

	fmt.Printf("\n=== Undetermined Propositions ===\n")

	for k, v := range invalidCache {
		_, found := cache[k]
		if found {
			continue
			//fmt.Printf(" [Dup %s, Actually Valid]", idx)
		}
		fmt.Printf("%-40s %s\n", v, k)
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
