package graph

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"

	svg "github.com/ajstarks/svgo"
)

type point struct {
	x float64
	y float64
}

func dist(p1 point, p2 point) float64 {
	return (p1.x-p2.x)*(p1.x-p2.x) + (p1.y-p2.y)*(p1.y-p2.y)
}

// Layout def
type Layout map[string]point

func drawGraph(canvas *svg.SVG, g Graph, scale int) {
	// get layout
	layout := geneticLayout(g, scale)

	for _, node := range g.nodes {
		p := layout[node]
		canvas.Circle(int(p.x), int(p.y), 8, "fill:red; stroke:red")
		text := fmt.Sprintf("%v", node)
		canvas.Text(int(p.x), int(p.y)-10, text, "text-anchor:middle;font-size:15px;fill:blue")
	}

	for node, adj := range g.adjList {
		st := layout[node]
		for _, neighb := range adj {
			fin := layout[neighb]
			//canvas.Arc(st.x, st.y, 5, 5, 0, false, false, fin.x, fin.y, "fill:none; stroke:blue")
			canvas.Line(int(st.x), int(st.y), int(fin.x), int(fin.y), "stroke:blue")
		}
	}
}

// arc layout
func arcLayout(g Graph, scale int) Layout {

	l := make(Layout)
	// use arc layout layout
	n := g.nodes
	dx := float64(scale / (len(n) + 1))
	x := dx
	y := float64(scale / 2)
	for _, node := range n {
		l[node] = point{x: x, y: y}
		x += dx
	}
	return l
}

func randomLayout(g Graph, scale int) Layout {
	l := make(Layout)
	n := g.nodes
	rand.Seed(time.Now().UnixNano())
	for _, node := range n {
		l[node] = point{x: float64(scale) * rand.Float64(), y: float64(scale) * rand.Float64()}
	}

	return l
}

func geneticLayout(g Graph, scale int) Layout {
	// reference: https://www.emis.de/journals/DM/v92/art5.pdf
	niter := 100
	rarray := []float64{}
	seeds := []Layout{}
	n := g.nodes
	rand.Seed(time.Now().UnixNano())

	// initial seeds
	for i := 0; i < 10; i++ {
		l := make(Layout)
		for _, node := range n {
			l[node] = point{x: float64(scale) * rand.Float64(), y: float64(scale) * rand.Float64()}
		}
		seeds = append(seeds, l)
	}

	//fmt.Println("seeds", seeds)

	// iteration loop
	for iter := 0; iter < niter; iter++ {

		// new generation from seeds
		larray := generateLayouts(seeds)
		// fmt.Println(iter, seeds)
		// rank current generation
		rarray = []float64{}
		for i := 0; i < len(larray); i++ {
			rarray = append(rarray, rank(larray[i], g.adjList))
		}
		fmt.Println(iter, rarray)

		// select seeds for next gen
		iarray := argmax(rarray)
		// fmt.Println(iter, iarray)
		seeds = []Layout{}
		for i := 0; i < len(iarray); i++ {
			seeds = append(seeds, larray[iarray[i]])
		}

	} // end iteration

	return seeds[0]
}

func argmax(rarray []float64) []int {
	iarray := []int{}
	pq := make(PriorityQueue, len(rarray))
	i := 0
	for value, priority := range rarray {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	for i := 0; i < 10; i++ {
		item := heap.Pop(&pq).(*Item)
		iarray = append(iarray, item.value)
	}
	return iarray
}

func generateLayouts(seeds []Layout) []Layout {
	larray := []Layout{}
	// copy seeds as is
	for i := 0; i < len(seeds); i++ {
		larray = append(larray, seeds[i])
	}

	// cross seeds two at a time
	for i := 0; i < len(seeds); i++ {
		for j := i + 1; j < len(seeds); j++ {
			larray = append(larray, cross(seeds[i], seeds[j]))
			larray = append(larray, cross(seeds[i], seeds[j]))
		}
	}

	return larray
}

func cross(s1 Layout, s2 Layout) Layout {
	l := make(Layout)
	rand.Seed(time.Now().UnixNano())
	i := 0
	for key, _ := range s1 {
		rnd := rand.Intn(2)
		if rnd == 0 {
			l[key] = s1[key]
		} else {
			l[key] = s2[key]
		}
		i++
	}

	return l
}

func rank(l Layout, adj adjacencyList) float64 {

	rep := 0.0
	att := 0.0

	for key, val := range l {
		for _, v := range adj[key] {
			att -= (0.001 * dist(l[v], val))
		}
	}

	for key, val := range l {
		for k, v := range l {
			if k != key {
				rep += (100 / (dist(v, val)))
			}
		}
	}

	return att - rep
}
