package CsvInput

import (
	"RemoveDuplicates"
	"os"
	"fmt"
	"encoding/csv"
	"bufio"
	"io"
	"log"
	"graph"
	"strconv"
	"strings"
)

type ArrayProcessedData struct {
	Array []ProcessedData
}
type ProcessedData struct {
	FromTo string
	Path   Path
	Value  []string
}

//Path type to save connection between Path and Id
type Path struct {
	Path1 []string
	Path2 []string
	Path3 []string
	Id1   []int
	Id2   []int
	Id3   []int
}

//type for get data from csv
func CsvInputFunc(name string) []RemoveDuplicates.Data {
	csvFile, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ';'
	var data []RemoveDuplicates.Data
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, RemoveDuplicates.Data{
			ID:         line[0],
			DepStation: line[1],
			ArrStation: line[2],
			Value:      line[3],
			DepTime:    line[4],
			ArrTime:    line[5],
		})

	}
	RemoveDuplicates.RemoveDuplicatesFromData(data)
	AlgorithmAccessor(data)
	return data
}
func AlgorithmAccessor(data []RemoveDuplicates.Data) {
	var temp []string
	var edges []string
	graph := graphs.NewGraph()
	temp = append(temp, data[0].DepStation)
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data); j++ {
			if data[i].DepStation == temp[len(temp)-1] {

			} else {
				//graph.AddVertex(data[i].DepStation, nil)
				temp = append(temp, data[i].DepStation)
			}
		}
	}
	edges = RemoveDuplicates.RemoveDuplicatesUnordered(temp)

	for i := 0; i < len(edges); i++ {
		graph.AddVertex(edges[i], nil)

	}
	// Value from string to float64
	for i := 0; i < len(data); i++ {
		b, err := strconv.ParseFloat(data[i].Value, 64)
		if err != nil {
			fmt.Println(err)
		}

		f, err := strconv.Atoi(data[i].ID)
		graph.AddEdge(data[i].DepStation, data[i].ArrStation, b, f, nil)
	}
	var ProcData []ProcessedData
	var GlobalValueDataHolder [][]string
	var GlobalPathDataHolder []Path
	var GlobalFromToDataHolder []string
	var counter1 = 0
	GlobalValueDataHolder = make([][]string, len(edges)*(len(edges)-1))
	GlobalPathDataHolder = make([]Path, len(edges)*(len(edges)-1))
	for i := 0; i < len(edges); i++ {
		for j := 0; j < len(edges); j++ {
			if edges[i] == edges[j] {

			} else {
				counter1++
				dist, path, _ := graph.Yen(edges[i], edges[j], 3)
				fmt.Println("The distance from: "+edges[i]+" to: "+edges[j]+" is ", dist)
				fmt.Println("The path from:"+edges[i]+" to: "+edges[j]+" is: ", path)
				GlobalFromToDataHolder = append(GlobalFromToDataHolder, "from: "+edges[i]+" to: "+edges[j])
				var IDValueHolder []graphs.Id
				var ValueHolder []string
				for l := 0; l < len(dist); l++ {
					var EachIDValueHolder graphs.Id
					var ToStringValueHolder string
					for k := 0; k < len(dist); k++ {
						EachIDValueHolder = dist[l]
						ToStringValueHolder = fmt.Sprintf("%v", EachIDValueHolder)
					}
					ValueHolder = append(ValueHolder, ToStringValueHolder)
					GlobalValueDataHolder[counter1-1] = ValueHolder
				}
				var IDPathHolder []graphs.Id
				var PathHolder Path
				for l := 0; l < len(path); l++ {
					var EachIDPathHolder []graphs.Id
					var ToStringPathHolder []string
					for k := 0; k < len(path[l]); k++ {
						EachIDPathHolder = path[l]
						ToStringPathHolder = append(ToStringPathHolder, fmt.Sprintf("%v", EachIDPathHolder[k]))
						if l == 0 {
							PathHolder.Path1 = ToStringPathHolder
						}
						if l == 1 {
							PathHolder.Path2 = ToStringPathHolder
						}
						if l == 2 {
							PathHolder.Path3 = ToStringPathHolder
						}
					}
					GlobalPathDataHolder[counter1-1] = PathHolder

				}
				for l := 0; l < len(dist); l++ {
					var a graphs.Id
					for k := 0; k < len(dist); k++ {
						a = dist[l]
					}
					IDValueHolder = append(IDValueHolder, a)
				}

				for l := 0; l < len(path); l++ {
					var a graphs.Id
					for k := 0; k < len(path[l]); k++ {
						a = path[l]
					}
					IDPathHolder = append(IDPathHolder, a)
				}
			}
		}
	}
	ProcData = make([]ProcessedData, counter1)
	for _ = range ProcData {
		for h, v := range GlobalValueDataHolder {
			ProcData[h].Value = v

		}
		/*for _, v := range ProcData[z].Value {
			fmt.Println(v)
		}*/
	}
	for _ = range ProcData {
		for h, v := range GlobalPathDataHolder {
			ProcData[h].Path = v
		}
		/*for _, v := range ProcData[z].Path.Path1 {
			fmt.Println(v)
		}*/
	}
	for z := range ProcData {
		ProcData[z].FromTo = GlobalFromToDataHolder[z]
	}
	/*for _, v := range ProcData[z].Path.Path1 {
		fmt.Println(v)
	}*/

	var NewData ArrayProcessedData
	NewData.Array = ProcData
	NewData.DataProcessor(graph)

	fmt.Println("You have:", len(edges), "stations")
	fmt.Println(edges)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your current station: ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	reader1 := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your destination station: ")
	text1, _ := reader1.ReadString('\n')
	text1 = strings.Replace(text1, "\n", "", -1)

	var TextOutPut string
	TextOutPut = "from: " + text + " to: " + text1
	for _, v := range NewData.Array {
		if TextOutPut == v.FromTo {
			fmt.Println("You have 3 path varieties: \n")

			//First path

			fmt.Println("First one is: ")
			for _, v := range v.Path.Path1 {
				fmt.Print("->")
				fmt.Print(v)
			}
			fmt.Println("\nTrains Id's: ")
			for _, v := range v.Path.Id1 {
				fmt.Print("->")
				fmt.Print(v)
			}
			//fmt.Println(v.Path.Id1)
			fmt.Println("\nThis path will cost for you: ")
			fmt.Println(v.Value[0], "\n")

			//Second path

			fmt.Println("Second one is: ")
			for _, v := range v.Path.Path2 {
				fmt.Print("->")
				fmt.Print(v)
			}
			fmt.Println("\nTrains Id's: ")
			for _, v := range v.Path.Id2 {
				fmt.Print("->")
				fmt.Print(v)
			}
			//fmt.Println(v.Path.Id1)
			fmt.Println("\nThis path will cost for you: ")
			fmt.Println(v.Value[1], "\n")

			//Third path

			fmt.Println("Third one is: ")
			for _, v := range v.Path.Path3 {
				fmt.Print("->")
				fmt.Print(v)
			}
			fmt.Println("\nTrains Id's: ")
			for _, v := range v.Path.Id3 {
				fmt.Print("->")
				fmt.Print(v)
			}
			//fmt.Println(v.Path.Id1)
			fmt.Println("\nThis path will cost for you: ")
			fmt.Println(v.Value[2])
			fmt.Println("Do you want to take more information about trains (y/n) ?")
			reader3 := bufio.NewReader(os.Stdin)
			text3, _ := reader3.ReadString('\n')
			text3 = strings.Replace(text3, "\n", "", -1)
			if text3 == "y" {
				fmt.Println("Information about first route: ")
				for _, v := range v.Path.Id1 {
					for i := 0; i < len(data); i++ {
						a, _ := strconv.Atoi(data[i].ID)

						if v == a {
							fmt.Println(data[i])
						}
					}

				}
				fmt.Println("Information about second route: ")
				for _, v := range v.Path.Id2 {
					for i := 0; i < len(data); i++ {
						a, _ := strconv.Atoi(data[i].ID)

						if v == a {
							fmt.Println(data[i])
						}
					}

				}
				fmt.Println("Information about third route: ")
				for _, k := range v.Path.Id3 {
					for i := 0; i < len(data); i++ {
						a, _ := strconv.Atoi(data[i].ID)

						if k == a {
							fmt.Println(data[i])
						}
					}

				}
			}
			if text3 == "n" {
			}
		}
	}
}
func (ProcData ArrayProcessedData) DataProcessor(graph *graphs.Graph) () {
	for i := 0; i < len(ProcData.Array); i++ {
		//var iml Path
		var id1, id2, id3 []string
		var index1, index2, index3 int
		index1, index2, index3 = -1, -1, -1
		for _, v := range ProcData.Array[i].Path.Path1 {
			id1 = append(id1, v)
			index1++
		}
		for _, v := range ProcData.Array[i].Path.Path2 {
			id2 = append(id2, v)
			index2++
		}
		for _, v := range ProcData.Array[i].Path.Path3 {
			id3 = append(id3, v)
			index3++
		}

		if len(id1) == 1 {
			fmt.Println("Error")
		}
		if len(id1) == 2 {

			ProcData.Array[i].Path.Id1 = append(ProcData.Array[i].Path.Id1, graph.Egress[id1[index1-1]][id1[index1]].Id)
			//fmt.Println("Маршрут 1: ", ProcData.Array[i].Path.Id1)
		}
		if len(id1) == 3 {
			ProcData.Array[i].Path.Id1 = append(ProcData.Array[i].Path.Id1, graph.Egress[id1[index1-2]][id1[index1-1]].Id)
			ProcData.Array[i].Path.Id1 = append(ProcData.Array[i].Path.Id1, graph.Egress[id1[index1-1]][id1[index1]].Id)
			//fmt.Println("Маршрут 1: ", ProcData.Array[i].Path.Id1)

		}
		if len(id1) == 4 {
			ProcData.Array[i].Path.Id1 = append(ProcData.Array[i].Path.Id1, graph.Egress[id1[index1-3]][id1[index1-2]].Id)
			ProcData.Array[i].Path.Id1 = append(ProcData.Array[i].Path.Id1, graph.Egress[id1[index1-2]][id1[index1-1]].Id)
			ProcData.Array[i].Path.Id1 = append(ProcData.Array[i].Path.Id1, graph.Egress[id1[index1-1]][id1[index1]].Id)
			//fmt.Println("Маршрут 1: ", ProcData.Array[i].Path.Id1)

		}
		if len(id1) == 5 {
			ProcData.Array[i].Path.Id1 = append(ProcData.Array[i].Path.Id1, graph.Egress[id1[index1-4]][id1[index1-3]].Id)
			ProcData.Array[i].Path.Id1 = append(ProcData.Array[i].Path.Id1, graph.Egress[id1[index1-3]][id1[index1-2]].Id)
			ProcData.Array[i].Path.Id1 = append(ProcData.Array[i].Path.Id1, graph.Egress[id1[index1-2]][id1[index1-1]].Id)
			ProcData.Array[i].Path.Id1 = append(ProcData.Array[i].Path.Id1, graph.Egress[id1[index1-1]][id1[index1]].Id)
			//fmt.Println("Маршрут 1: ", ProcData.Array[i].Path.Id1)
		}

		//fmt.Println(id1)

		if len(id2) == 1 {
			fmt.Println("Error")
		}
		if len(id2) == 2 {

			ProcData.Array[i].Path.Id2 = append(ProcData.Array[i].Path.Id2, graph.Egress[id2[index2-1]][id2[index2]].Id)
			//fmt.Println("Маршрут 2: ", ProcData.Array[i].Path.Id2)
		}
		if len(id2) == 3 {
			ProcData.Array[i].Path.Id2 = append(ProcData.Array[i].Path.Id2, graph.Egress[id2[index2-2]][id2[index2-1]].Id)
			ProcData.Array[i].Path.Id2 = append(ProcData.Array[i].Path.Id2, graph.Egress[id2[index2-1]][id2[index2]].Id)
			//fmt.Println("Маршрут 2: ", ProcData.Array[i].Path.Id2)

		}
		if len(id2) == 4 {
			ProcData.Array[i].Path.Id2 = append(ProcData.Array[i].Path.Id2, graph.Egress[id2[index2-3]][id2[index2-2]].Id)
			ProcData.Array[i].Path.Id2 = append(ProcData.Array[i].Path.Id2, graph.Egress[id2[index2-2]][id2[index2-1]].Id)
			ProcData.Array[i].Path.Id2 = append(ProcData.Array[i].Path.Id2, graph.Egress[id2[index2-1]][id2[index2]].Id)
			//fmt.Println("Маршрут 2: ", ProcData.Array[i].Path.Id2)

		}
		if len(id2) == 5 {
			ProcData.Array[i].Path.Id2 = append(ProcData.Array[i].Path.Id2, graph.Egress[id2[index2-4]][id2[index2-3]].Id)
			ProcData.Array[i].Path.Id2 = append(ProcData.Array[i].Path.Id2, graph.Egress[id2[index2-3]][id2[index2-2]].Id)
			ProcData.Array[i].Path.Id2 = append(ProcData.Array[i].Path.Id2, graph.Egress[id2[index2-2]][id2[index2-1]].Id)
			ProcData.Array[i].Path.Id2 = append(ProcData.Array[i].Path.Id2, graph.Egress[id2[index2-1]][id2[index2]].Id)
			//fmt.Println("Маршрут 2: ", ProcData.Array[i].Path.Id2)
		}
		//fmt.Println(id2)

		if len(id3) == 1 {
			fmt.Println("Error")
		}
		if len(id3) == 2 {

			ProcData.Array[i].Path.Id3 = append(ProcData.Array[i].Path.Id3, graph.Egress[id3[index3-1]][id3[index3]].Id)
			//fmt.Println("Маршрут 3: ", ProcData.Array[i].Path.Id3)
		}
		if len(id3) == 3 {
			ProcData.Array[i].Path.Id3 = append(ProcData.Array[i].Path.Id3, graph.Egress[id3[index3-2]][id3[index3-1]].Id)
			ProcData.Array[i].Path.Id3 = append(ProcData.Array[i].Path.Id3, graph.Egress[id3[index3-1]][id3[index3]].Id)
			//fmt.Println("Маршрут 3: ", ProcData.Array[i].Path.Id3)

		}
		if len(id3) == 4 {
			ProcData.Array[i].Path.Id3 = append(ProcData.Array[i].Path.Id3, graph.Egress[id3[index3-3]][id3[index3-2]].Id)
			ProcData.Array[i].Path.Id3 = append(ProcData.Array[i].Path.Id3, graph.Egress[id3[index3-2]][id3[index3-1]].Id)
			ProcData.Array[i].Path.Id3 = append(ProcData.Array[i].Path.Id3, graph.Egress[id3[index3-1]][id3[index3]].Id)
			//fmt.Println("Маршрут 3: ", ProcData.Array[i].Path.Id3)

		}
		if len(id3) == 5 {
			ProcData.Array[i].Path.Id3 = append(ProcData.Array[i].Path.Id3, graph.Egress[id3[index3-4]][id3[index3-3]].Id)
			ProcData.Array[i].Path.Id3 = append(ProcData.Array[i].Path.Id3, graph.Egress[id3[index3-3]][id3[index3-2]].Id)
			ProcData.Array[i].Path.Id3 = append(ProcData.Array[i].Path.Id3, graph.Egress[id3[index3-2]][id3[index3-1]].Id)
			ProcData.Array[i].Path.Id3 = append(ProcData.Array[i].Path.Id3, graph.Egress[id3[index3-1]][id3[index3]].Id)
			//fmt.Println("Маршрут 3: ", ProcData.Array[i].Path.Id3)
		}
		//fmt.Println(id3)

	}

}