package graphs

import (
	"fmt"
	"math"
	"github.com/starwander/GoFibonacciHeap"
)



// Id uniquely identify a vertex.
type Id interface{}

// Graph is made up of vertices and edges.
// Vertices in the graph must have an unique id.
// Each edges in the graph connects two vertices directed with a weight.
type Graph struct {
	vertices map[Id]*vertex
	Egress   map[Id]map[Id]*edge
	ingress  map[Id]map[Id]*edge
}

type vertex struct {
	self   interface{}
	enable bool
}

type edge struct {
	self    interface{}
	weight  float64
	enable  bool
	changed bool
	Id      int
}
// NewGraph creates a new empty graph.
func NewGraph() *Graph {
	graph := new(Graph)
	graph.vertices = make(map[Id]*vertex)
	graph.Egress = make(map[Id]map[Id]*edge)
	graph.ingress = make(map[Id]map[Id]*edge)

	return graph
}
func (edge *edge) getId() int {
	return edge.Id
}
func (edge *edge) getWeight() float64 {
	return edge.weight
}

// GetVertex get a vertex by input id.
// Try to get a vertex not in the graph will get an error.
func (graph *Graph) GetVertex(id Id) (vertex interface{}, err error) {
	if v, exists := graph.vertices[id]; exists {
		vertex = v.self
		return
	}

	err = fmt.Errorf("Vertex %v is not found", id)
	return
}

// GetEdge gets the edge between the two vertices by input ids.
// Try to get the edge from or to a vertex not in the graph will get an error.
// Try to get the edge between two disconnected vertices will get an error.
func (graph *Graph) GetEdge(from Id, to Id) (interface{}, error) {
	if _, exists := graph.vertices[from]; !exists {
		return nil, fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := graph.vertices[to]; !exists {
		return nil, fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if edge, exists := graph.Egress[from][to]; exists {
		return edge.self, nil
	}

	return nil, fmt.Errorf("Edge from %v to %v is not found", from, to)
}

// GetEdgeWeight gets the weight of the edge between the two vertices by input ids.
// Try to get the weight of the edge from or to a vertex not in the graph will get an error.
// Try to get the weight of the edge between two disconnected vertices will get +Inf.
func (graph *Graph) GetEdgeWeight(from Id, to Id) (float64, error) {
	if _, exists := graph.vertices[from]; !exists {
		return math.Inf(1), fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := graph.vertices[to]; !exists {
		return math.Inf(1), fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if edge, exists := graph.Egress[from][to]; exists {
		return edge.weight, nil
	}

	return math.Inf(1), nil
}

// AddVertex adds a new vertex into the graph.
// Try to add a duplicate vertex will get an error.
func (graph *Graph) AddVertex(id Id, v interface{}) error {
	if _, exists := graph.vertices[id]; exists {
		return fmt.Errorf("Vertex %v is duplicate", id)
	}

	graph.vertices[id] = &vertex{v, true}
	graph.Egress[id] = make(map[Id]*edge)
	graph.ingress[id] = make(map[Id]*edge)

	return nil
}

// AddEdgeadds a new edge between the vertices by the input ids.
// Try to add an edge with -Inf weight will get an error.
// Try to add an edge from or to a vertex not in the graph will get an error.
// Try to add a duplicate edge will get an error.
func (graph *Graph) AddEdge(from Id, to Id, weight float64, identificator int, e interface{}) error {
	if weight == math.Inf(-1) {
		return fmt.Errorf("-inf weight is reserved for internal usage")
	}

	if _, exists := graph.vertices[from]; !exists {
		return fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := graph.vertices[to]; !exists {
		return fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if _, exists := graph.Egress[from][to]; exists {
		for i := 0; i < len(graph.Egress); i++ {
			//Check for the best edge
			if weight < graph.Egress[from][to].weight {
				graph.Egress[from][to] = &edge{e, weight, true, false, identificator}
				graph.ingress[to][from] = graph.Egress[from][to]
			}
		}
		return fmt.Errorf("Edge from %v to %v is duplicate", from, to)
	}

	graph.Egress[from][to] = &edge{e, weight, true, false, identificator}
	graph.ingress[to][from] = graph.Egress[from][to]

	return nil
}

// UpdateEdgeWeight updates the weight of the edge between vertices by the input ids.
// Try to update an edge with -Inf weight will get an error.
// Try to update an edge from or to a vertex not in the graph will get an error.
// Try to update an edge between disconnected vertices will get an error.
func (graph *Graph) UpdateEdgeWeight(from Id, to Id, weight float64) error {
	if weight == math.Inf(-1) {
		return fmt.Errorf("-inf weight is reserved for internal usage")
	}

	if _, exists := graph.vertices[from]; !exists {
		return fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := graph.vertices[to]; !exists {
		return fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if edge, exists := graph.Egress[from][to]; exists {
		edge.weight = weight
		return nil
	}

	return fmt.Errorf("Edge from %v to %v is not found", from, to)
}

// DeleteVertex deletes a vertex from the graph and gets the value of the vertex.
// Try to delete a vertex not in the graph will get an nil.
func (graph *Graph) DeleteVertex(id Id) interface{} {
	if vertex, exists := graph.vertices[id]; exists {
		for to := range graph.Egress[id] {
			delete(graph.ingress[to], id)
		}
		for from := range graph.ingress[id] {
			delete(graph.Egress[from], id)
		}
		delete(graph.Egress, id)
		delete(graph.ingress, id)
		delete(graph.vertices, id)

		return vertex.self
	}

	return nil
}

// DeleteEdge deletes the edge between the vertices by the input id from the graph and gets the value of edge.
// Try to delete an edge from or to a vertex not in the graph will get an error.
// Try to delete an edge between disconnected vertices will get a nil.
func (graph *Graph) DeleteEdge(from Id, to Id) interface{} {
	if _, exists := graph.vertices[from]; !exists {
		return nil
	}

	if _, exists := graph.vertices[to]; !exists {
		return nil
	}

	if edge, exists := graph.Egress[from][to]; exists {
		delete(graph.Egress[from], to)
		delete(graph.ingress[to], from)
		return edge.self
	}

	return nil
}

// AddVertexWithEdges adds a vertex value which implements Vertex interface.
// AddVertexWithEdges adds edges connected to the vertex at the same time, due to the Vertex interface can get the Edges.
// CheckIntegrity checks if any edge connects to or from unknown vertex.
// If the graph is integrate, nil is returned. Otherwise an error is returned.
func (graph *Graph) CheckIntegrity() error {
	for from, out := range graph.Egress {
		if _, exists := graph.vertices[from]; !exists {
			return fmt.Errorf("Vertex %v is not found", from)
		}
		for to := range out {
			if _, exists := graph.vertices[to]; !exists {
				return fmt.Errorf("Vertex %v is not found", to)
			}
		}
	}

	for to, in := range graph.ingress {
		if _, exists := graph.vertices[to]; !exists {
			return fmt.Errorf("Vertex %v is not found", to)
		}
		for from := range in {
			if _, exists := graph.vertices[from]; !exists {
				return fmt.Errorf("Vertex %v is not found", from)
			}
		}
	}

	return nil
}

// GetPathWeight gets the total weight along the path by input ids.
// It will get -Inf if the input path is nil or empty.
// It will get -Inf if the path contains vertex not in the graph.
// It will get +Inf if the path contains vertices not connected.
func (graph *Graph) GetPathWeight(path []Id) (totalWeight float64, totalId []int) {
	if len(path) == 0 {
		return math.Inf(-1), nil
	}

	if _, exists := graph.vertices[path[0]]; !exists {
		return math.Inf(-1), nil
	}

	for i := 0; i < len(path)-1; i++ {
		if _, exists := graph.vertices[path[i+1]]; !exists {
			return math.Inf(-1), nil
		}
		if edge, exists := graph.Egress[path[i]][path[i+1]]; exists {
			totalWeight += edge.getWeight()

			totalId = append(totalId, edge.Id)

		} else {
			return math.Inf(1), nil
		}
	}

	return totalWeight, totalId
}
func (graph *Graph) GetPathId(path []Id) (totalId int) {
	if len(path) == 0 {
		return
	}

	if _, exists := graph.vertices[path[0]]; !exists {
		return
	}

	for i := 0; i < len(path)-1; i++ {
		if _, exists := graph.vertices[path[i+1]]; !exists {
			return
		}
		if _, exists := graph.Egress[path[i]][path[i+1]]; exists {

			return totalId
		} else {
			return
		}
	}

	return
}

// DisableEdge disables the edge for further calculation.
func (graph *Graph) DisableEdge(from, to Id) {
	graph.Egress[from][to].enable = false
}

// DisableVertex disables the vertex for further calculation.
func (graph *Graph) DisableVertex(vertex Id) {
	for _, edge := range graph.Egress[vertex] {
		edge.enable = false
	}
}

// DisablePath disables all the vertices in the path for further calculation.
func (graph *Graph) DisablePath(path []Id) {
	for _, vertex := range path {
		graph.DisableVertex(vertex)
	}
}

// Reset enables all vertices and edges for further calculation.
func (graph *Graph) Reset() {
	for _, out := range graph.Egress {
		for _, edge := range out {
			edge.enable = true
		}
	}
}

// Yen gets top k shortest loopless path between two vertex in the graph.
func (graph *Graph) Yen(source, destination Id, topK int) ([]float64, [][]Id, error) {
	var err error
	var i, j, k int
	var spurWeight float64
	var dijkstraDist map[Id]float64
	var dijkstraPrev map[Id]Id

	distTopK := make([]float64, topK)
	pathTopK := make([][]Id, topK)

	for i := 0; i < topK; i++ {
		distTopK[i] = math.Inf(1)
	}

	dijkstraDist, dijkstraPrev, _, err = graph.Dijkstra(source)
	if err != nil {
		return nil, nil, err
	}
	distTopK[0] = dijkstraDist[destination]
	pathTopK[0] = getPath(dijkstraPrev, destination)
	for k = 1; k < topK; k++ {
		for i = 0; i < len(pathTopK[k-1])-1; i++ {
			for j = 0; j < k; j++ {
				if isShareRootPath(pathTopK[j], pathTopK[k-1][:i+1]) {
					graph.DisableEdge(pathTopK[j][i], pathTopK[j][i+1])

				}
			}
			graph.DisablePath(pathTopK[k-1][:i])

			dijkstraDist, dijkstraPrev, _, err = graph.Dijkstra(pathTopK[k-1][i])

			spurWeight, _ = graph.GetPathWeight(pathTopK[k-1][:i+1])
			spurWeight += dijkstraDist[destination]
			if spurWeight < distTopK[k] {
				distTopK[k] = spurWeight
				pathTopK[k] = mergePath(pathTopK[k-1][:i], getPath(dijkstraPrev, destination))

			}
			graph.Reset()
		}
	}

	return distTopK, pathTopK, nil
}

func isShareRootPath(path, rootPath []Id) bool {
	if len(path) < len(rootPath) {
		return false
	}

	for i := 0; i < len(rootPath); i++ {
		if path[i] != rootPath[i] {
			return false
		}
	}

	return true
}

func mergePath(path1, path2 []Id) []Id {
	newPath := []Id{}
	newPath = append(newPath, path1...)
	newPath = append(newPath, path2...)

	return newPath
}

// Dijkstra gets the shortest path from one vertex to all other vertices in the graph.
func (graph *Graph) Dijkstra(source Id) (dist map[Id]float64, prev map[Id]Id, id []int, err error) {
	if _, exists := graph.vertices[source]; !exists {
		return nil, nil, nil, fmt.Errorf("Vertex %v is not existed", source)
	}

	dist = make(map[Id]float64)
	prev = make(map[Id]Id)
	heap := fibHeap.NewFibHeap()

	for id := range graph.vertices {
		prev[id] = nil
		if id != source {
			dist[id] = math.Inf(1)
			heap.Insert(id, math.Inf(1))
		} else {
			dist[id] = 0
			heap.Insert(id, 0)
		}
	}

	for heap.Num() != 0 {
		min, _ := heap.ExtractMin()
		for to, edge := range graph.Egress[min] {
			if edge.getWeight() < 0 {
				return nil, nil, nil, fmt.Errorf("Negative weight form vertex %v to vertex %v is not allowed", min, to)
			}
			if !edge.enable {
				continue
			}
			if dist[min]+edge.getWeight() < dist[to] {
				heap.DecreaseKey(to, dist[min]+edge.getWeight())
				prev[to] = min
				dist[to] = dist[min] + edge.getWeight()
			}

		}

	}

	return
}

//[1179 1146 1350 909 1398 957 1148 1155 1492 1068 1013 1140 1125 1342 1137 1177 919 1139 920 1126 1156 1349 1365 1171]
func getPath(prev map[Id]Id, lastNode Id) (path []Id) {
	prevNode := prev[lastNode]
	if prevNode == nil {
		return nil
	}

	reversePath := []Id{lastNode}
	for ; prevNode != nil; prevNode = prev[prevNode] {
		reversePath = append(reversePath, prevNode)
	}

	path = make([]Id, len(reversePath))
	for index, node := range reversePath {
		path[len(reversePath)-index-1] = node
	}

	return
}