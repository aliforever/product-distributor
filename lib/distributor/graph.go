package distributor

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/multi"
	"gonum.org/v1/gonum/graph/path"
	"math"
	"sort"
)

// RequiredPacks represents a mapping between the size of a pack and the number
// of packs required.
type RequiredPacks map[int]int

// GraphPackCalculator is responsible for calculating the number of packs needed
// given a set of available pack sizes.
type GraphPackCalculator struct {
	PackSizes []int
}

// quantityGraph represents a directed multigraph with weighted edges. It keeps
// track of quantities and corresponding packs for those quantities.
type quantityGraph struct {
	packSizeCount int
	candidates    map[int]quantityNode
	*multi.WeightedDirectedGraph
}

// quantityNode represents a node within the quantityGraph. Each node has an
// associated quantity which may or may not be a target quantity.
type quantityNode struct {
	quantity int
}

// headroomMultiplier is a constant multiplier to give room for permutations.
const headroomMultiplier int = 50

// Calculate determines the required number of packs for the given quantity. It
// utilizes a graph-based approach to generate permutations and derive an optimal solution.
func (c GraphPackCalculator) Calculate(quantity int) (RequiredPacks, error) {
	packs := make(RequiredPacks)
	if quantity <= 0 {
		return packs, nil
	}

	c.initializePacksBySize(&quantity, packs)
	qGraph := c.createInitialGraph(quantity)
	c.generatePermutations(qGraph, quantity) // Here, pass the quantity directly
	qGraph.trimUnnecessaryNodes()

	packsFromShortestPath := c.derivePacksFromShortestPath(qGraph, quantity) // Also pass the quantity here
	mergePackMaps(packs, packsFromShortestPath)

	return packs, nil
}

// Weight fetches the uniform cost (weight) between two nodes in the quantityGraph.
func (g quantityGraph) Weight(xid, yid int64) (w float64, ok bool) {
	return path.UniformCost(g)(xid, yid)
}

// ID returns the unique ID for a quantityNode, which is its quantity value.
func (n quantityNode) ID() int64 {
	return int64(n.quantity)
}

// trimUnnecessaryNodes removes nodes that aren't crucial for determining the shortest path.
func (g *quantityGraph) trimUnnecessaryNodes() {
	candidateNode := g.closestCandidate()

	// remove other candidates from the graph
	for _, node := range g.candidates {
		if node != candidateNode {
			g.RemoveNode(node.ID())
		}
	}

	// remove nodes which don't have any edges going out
	var retraverse bool
	for {
		retraverse = false
		it := g.Nodes()
		for it.Next() {
			if node := it.Node(); node != candidateNode && len(graph.NodesOf(g.From(node.ID()))) == 0 {
				g.RemoveNode(node.ID())
				retraverse = true
			}
		}
		if !retraverse {
			break
		}
	}
}

// derivePacksFromShortestPath identifies the required packs by analyzing the shortest path
// in the graph between the starting quantity and the closest candidate.
func (c GraphPackCalculator) derivePacksFromShortestPath(qGraph *quantityGraph, quantity int) RequiredPacks {
	resultPacks := make(RequiredPacks)
	candidateNode := qGraph.closestCandidate()

	shortest, _ := path.AStar(quantityNode{quantity}, candidateNode, qGraph, nil)
	path, _ := shortest.To(candidateNode.ID())
	pathLength := len(path)

	for i, currentNode := range path {
		nextIndex := i + 1
		if nextIndex >= pathLength {
			break
		}

		lines := qGraph.WeightedLines(currentNode.ID(), path[nextIndex].ID())
		lines.Next()
		resultPacks[int(lines.WeightedLine().Weight())]++
	}
	return resultPacks
}

// generatePermutations produces permutations of available pack sizes in the graph to reach the target quantity.
func (c GraphPackCalculator) generatePermutations(qGraph *quantityGraph, quantity int) { // Added quantity parameter
	sizes := c.PackSizes
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	for i := len(sizes); i >= 1; i-- {
		availableSizes := make([]int, i)
		copy(availableSizes, sizes[:i])
		qGraph.subtractPacks(quantityNode{quantity}, availableSizes)
	}
}

// createInitialGraph initializes a quantityGraph with the provided quantity as its root node.
func (c GraphPackCalculator) createInitialGraph(quantity int) *quantityGraph {
	qGraph := quantityGraph{
		packSizeCount:         len(c.PackSizes),
		candidates:            make(map[int]quantityNode),
		WeightedDirectedGraph: multi.NewWeightedDirectedGraph(),
	}
	rootNode := quantityNode{quantity}
	qGraph.AddNode(rootNode)
	return &qGraph
}

// initializePacksBySize adjusts the initial number of largest packs needed to approach the target quantity.
func (c GraphPackCalculator) initializePacksBySize(quantity *int, packs RequiredPacks) {
	sizes := c.PackSizes
	sort.Ints(sizes)

	permutationClamp := sum(sizes) * headroomMultiplier
	if *quantity > permutationClamp {
		largestSize := sizes[len(sizes)-1]
		packs[largestSize] = int(math.Floor(float64(*quantity-permutationClamp) / float64(largestSize)))
		*quantity -= packs[largestSize] * largestSize
	}
}

// subtractPacks subtracts available pack sizes from the current quantity and adds resulting nodes
// and edges to the graph, recursively.
func (g *quantityGraph) subtractPacks(n quantityNode, packSizes []int) {
	// stop generating permutations if we've found more paths to 0 than available pack sizes
	if nodesToZero := g.To(int64(0)); nodesToZero.Len() >= g.packSizeCount {
		return
	}

	for _, size := range packSizes {
		// find or create a node by the subtracted quantity
		nextQuantity := n.quantity - size
		nextNode := quantityNode{nextQuantity}
		if existingNode := g.Node(nextNode.ID()); existingNode == nil {
			g.AddNode(nextNode)
		}

		// maintain unique weights for edges between two quantities to avoid unnecessary recalculations
		weight := float64(size)
		if g.hasWeightedLineFromTo(n, nextNode, weight) {
			continue
		}

		// link the nodes by pack size
		g.SetWeightedLine(g.NewWeightedLine(n, nextNode, weight))

		// track nodes which satisfy the required quantity, stopping at this depth
		if nextQuantity <= 0 {
			g.candidates[nextQuantity] = nextNode
			continue
		}

		// subtract from the next quantity, increasing depth
		g.subtractPacks(nextNode, packSizes)
	}
}

// hasWeightedLineFromTo checks if there's already an edge with a specific weight between two nodes.
func (g *quantityGraph) hasWeightedLineFromTo(from graph.Node, to graph.Node, weight float64) bool {
	for _, line := range graph.WeightedLinesOf(g.WeightedLines(from.ID(), to.ID())) {
		if line.Weight() == weight {
			return true
		}
	}
	return false
}

// closestCandidate identifies the node that's closest to zero from the candidates.
func (g *quantityGraph) closestCandidate() quantityNode {
	// create a slice of quantities from the map keys
	quantities := make([]int, len(g.candidates))
	i := 0
	for k := range g.candidates {
		quantities[i] = k
		i++
	}

	// reverse sort so the closest candidate is first
	sort.Sort(sort.Reverse(sort.IntSlice(quantities)))

	return g.candidates[quantities[0]]
}

// pruneNodes is an older version of trimUnnecessaryNodes. It might be deprecated or unused.
func (g *quantityGraph) pruneNodes(candidate graph.Node) {
	// remove other candidates from the graph
	for _, node := range g.candidates {
		if node != candidate {
			g.RemoveNode(node.ID())
		}
	}

	// remove nodes which don't have any edges going out
	var retraverse bool
	for {
		retraverse = false
		it := g.Nodes()
		for it.Next() {
			if node := it.Node(); node != candidate && len(graph.NodesOf(g.From(node.ID()))) == 0 {
				g.RemoveNode(node.ID())
				retraverse = true
			}
		}
		if !retraverse {
			break
		}
	}
}

// sum calculates the sum of the integers in the provided slice.
func sum(array []int) int {
	sum := 0
	for _, v := range array {
		sum += v
	}
	return sum
}

// mergePackMaps merges the quantities of packs from one map into another.
func mergePackMaps(mainPacks, additionalPacks RequiredPacks) {
	for packSize, count := range additionalPacks {
		mainPacks[packSize] += count
	}
}
