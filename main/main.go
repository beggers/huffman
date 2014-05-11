// Ben Eggers
// GNU GPL'd

package main

import (
	"github.com/BenedictEggers/huffman"
	"os"
	"fmt"
)

// This will be the actual program to use the huffman tree to encode and decode
// files. It will accept a command-line argument giving it a file name, and 
// flags to tell it what to with that file (encode, or decode).

func main() {
	args := os.Args
	if len(args) != 4 {
		Usage(args[0])
	}
	t, err := huffman.MakeTreeFromText("filler")
	if err == nil {
		fmt.Println("Improperly-formatted text file.")
		os.Exit(1)
	}
	

}

// Prints out the usage of the program, and exits with an error status. Called
// if the user puts in bad command-line args.
func Usage(progName string) {
	fmt.Println("Usage:", progName, "-[e | d] FROM_FILE TO_FILE")
	os.Exit(1)
}