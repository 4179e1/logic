package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const nonPrefix = "non-"

var depthFlag int

func init() {

	flag.IntVar(&depthFlag, "d", 0, "Depth")
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

func A(s, p string, depthStr string) {
	str := fmt.Sprintf("All %s is %s", s, p)
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

func E(s, p string, depthStr string) {
	str := fmt.Sprintf("No %s is %s", s, p)
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

func I(s, p string, depthStr string) {

	str := fmt.Sprintf("Some %s is %s", s, p)
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

func O(s, p string, depthStr string) {
	str := fmt.Sprintf("Some %s is not %s", s, p)
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
