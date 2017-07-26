package bike

import (
	"fmt"
	"math"

	"github.com/oakmound/oak/alg/intgeom"
)

func ConnectGraph(nodes []intgeom.Point) [][]int {
	connections := make([][]int, len(nodes))
	// Prim's MST
	// We'll act like everything is adjacent to begin with
	mstNodes := make([]intgeom.Point, 1)
	mstNodes[0] = nodes[0]

	weights := make([]float64, len(nodes))
	for i := 0; i < len(weights); i++ {
		weights[i] = nodes[i].Distance(nodes[0])
	}
	weights[0] = math.MaxFloat64

	// Indicates from which index a node is being connected
	weightConnector := make([]int, len(nodes))
	for len(mstNodes) < len(nodes) {
		minIndex := 0
		minVal := math.MaxFloat64
		for i, w := range weights {
			if w < minVal {
				minVal = w
				minIndex = i
			}
		}
		mstNodes = append(mstNodes, nodes[minIndex])
		// Connect whoever represents the minimum weight to this
		connections[weightConnector[minIndex]] = append(connections[weightConnector[minIndex]], minIndex)
		connections[minIndex] = append(connections[minIndex], weightConnector[minIndex])
		// So we don't choose this again
		weights[minIndex] = math.MaxFloat64
	weightUpdate:
		for i, n := range nodes {
			// check n is not in mstNodes
			for _, n2 := range mstNodes {
				if n2 == n {
					continue weightUpdate
				}
			}
			dist := n.Distance(nodes[minIndex])
			if dist < weights[i] {
				weights[i] = dist
				weightConnector[i] = minIndex
			}
		}
	}
	fmt.Println(connections)
	return connections
}
