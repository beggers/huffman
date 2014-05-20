// Ben Eggers
// GNU GPL'd

package huffTree

import (
	"container/heap"
	"errors"
	"github.com/BenedictEggers/bitIO"
	"io/ioutil"
)

// The actual Huffman Tree and all associated functions. Will build up a
// tree from either a file defining the text file to be encoded or a file
// defining the Huffman Tree (see makeTreeFromText(), makeTreeFromTreeFile()) and
// will write a tree out to disk (see tree.writeToFile()). Uses the node struct found
// in "./node.go".

// HuffTree is really just a pointer to the root node of the tree.
type HuffTree struct {
	root *huffNode
}

// decode turns the bytes in fromFile into bytes in toFile, decompressed under
// the tree it is called on. On success, returns a nil error and returns a
// non-nil error otherwise. If fromFile exists before the call, it is deleted
// and replaced with the decompressed file.
func (t *HuffTree) DecodeText(fromFile, toFile string) (err error) {
	// Set up a BitReader on the file to decodes
	reader, err := bitIO.NewReader(fromFile)
	if err != nil {
		return err
	}

	// And open up our file to write to
	outFile, err := os.Create(toFile)
	if err != nil {
		return err
	}

	// Decode our bits, writing them out to disk every 1,000 bytes so as not
	// to use up all of main memory
	toWrite := make([]byte, 1000)
	currentIndex := 0
	bit, err := reader.ReadBit() // get the first bit

	// until we reach the end of the file...
	for err != nil {

		// Read out 1,000 bytes' worth of compressed data
		for currentIndex < len(toWrite) {
			// Process the current bit
			

			// And read a new bit
			bit, err = reader.ReadBit()
			if err != nil {
				return err
			}
		}

		// Write our bytes
		_, err = outFile.Write(toWrite)
		if err != nil {
			return err
		}

		// Re-zero our slice (so the last write won't write extra stuff)
		toWrite = make([]byte, len(toWrite))

		// And try to read another bit, so the loop will exit if we're at the EOF
		bit, err = reader.ReadBit()
	}

	// Great, we're done
	return nil
}

// encode turns the bytes in fromFile into bytes in toFile, compressed under
// the tree it is called on. On success, returns a nil error and returns a
// non-nil error otherwise.
func (t *HuffTree) EncodeText(fromFile, toFile string) (err error) {
	return errors.New("Undefined method")
}

// Write the tree out to a file at a file described by the passed string.
// Necessary to keep around if you plan on decoding files :)
func (t *HuffTree) WriteToFile(filename string) (err error) {
	return errors.New("Undefined method")
}

////////////////////////////////////////////////////////////////////////////////
//               Functions to make the tree.
////////////////////////////////////////////////////////////////////////////////

// makeTreeFromText takes in a text file and turns it into a HuffTree, which
// it then returns.
func MakeTreeFromText(filename string) (t *HuffTree, err error) {
	// Read the text byte-by-byte, building up a map of byte counts
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Scan the byte slice "buf" and count how many times each byte shows up
	counts := map[byte]uint32{}
	for _, elem := range buf {
		counts[elem] += 1
	}

	// Turn the counts into huffNodes
	nodes := make([]*huffNode, 0)
	for currentByte, byteCount := range counts {
		node := &huffNode{char: currentByte, count: byteCount}
		nodes = append(nodes, node)
	}

	// And make a tree
	return makeTreeFromNodeSlice(nodes)
}

// makeTreeFromTreeFile takes in a filname of a file in the same format TREE.writeToFile()
// puts out, and remakes a HuffTree from it.
func MakeTreeFromTreeFile(filename string) (t *HuffTree, err error) {
	return &HuffTree{}, errors.New("Undefined method")
}

// makeTreeFromNodeSlice makes a huffman tree from the passed slice of huffNodes.
// If len(nodes) == 0, returns a nil tree.
func makeTreeFromNodeSlice(nodes []*huffNode) (t *HuffTree, err error) {
	if len(nodes) == 0 {
		return nil, errors.New("Invalid node slice.")
	}

	// We're going to put the nodes in a heap, with low-ness determined
	// by the nodes' counts.
	nh := &nodeHeap{}
	heap.Init(nh)
	for _, node := range nodes {
		heap.Push(nh, node)
	}

	// Now, we're going to do the following:
	//
	// Until there's only one node in the heap:
	// 		Remove the lowest-count two nodes
	// 		Make a new node with those two as children, whose count is the
	//			sum of its childrens' counts
	//		Add that new node to the heap
	//
	// This will create an optimally-balanced tree, based on byte counts. For
	// more information, see http://en.wikipedia.org/wiki/Huffman_coding.
	for nh.Len() > 1 {
		nodeOne := heap.Pop(nh).(*huffNode)
		nodeTwo := heap.Pop(nh).(*huffNode)
		newNode := &huffNode{char: 0,
			count: nodeOne.count + nodeTwo.count,
			left:  nodeOne,
			right: nodeTwo}
		heap.Push(nh, newNode)
	}

	// Great, now there's only one node and it's the root of the tree!
	return &HuffTree{heap.Pop(nh).(*huffNode)}, nil
}
