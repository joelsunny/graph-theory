package graph

import (
	"os"
	"testing"

	svg "github.com/ajstarks/svgo"
)

// TestAdjacencyList tests initialization of graph from an adjacency list
func TestAdjacencyList(t *testing.T) {
	a := make(adjacencyList)
	a["a"] = []string{"b", "c"}
	a["c"] = []string{"d"}
	a["b"] = []string{"d"}
	a["d"] = []string{"a"}

	n := []string{"a", "b", "c", "d"}
	g := NewGraph(a, n)

	width := 512
	height := 512
	fo, _ := os.Create("out.svg")
	defer fo.Close()
	canvas := svg.New(fo)
	canvas.Start(width, height)
	drawGraph(canvas, g, 512)
	canvas.End()
}
