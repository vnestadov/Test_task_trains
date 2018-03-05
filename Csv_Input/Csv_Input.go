package Csv_Input

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vnestadov/Test_task_trains/Remove_Duplicates"
	"github.com/vnestadov/Test_task_trains/graph"
)

type InfoArray struct {
	Array []ProcessedData
}
type ProcessedData struct {
	FromTo string
	Path   RouteInfo
	Value  []string
}
type id_handler struct {
	id1, id2, id3          []string
	index1, index2, index3 int
}

//RouteInfo type to save connection between RouteInfo and Id
type RouteInfo struct {
	Route1 []string
	Route2 []string
	Route3 []string
	ID1    []int
	ID2    []int
	ID3    []int
}

func Input(name string) []Remove_Duplicates.Data {
	csvFile, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ';'
	var data []Remove_Duplicates.Data
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Remove_Duplicates.Data{
			ID:         line[0],
			DepStation: line[1],
			ArrStation: line[2],
			Value:      line[3],
			DepTime:    line[4],
			ArrTime:    line[5],
		})

	}
	Remove_Duplicates.Delete_Data(data)
	Algorithm_Starter(data)
	return data
}
func AddEdges(data []Remove_Duplicates.Data) (edges []string) {
	var temp []string
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
	edges = Remove_Duplicates.Delete(temp)
	return
}
func GetPathId(dist []float64, path [][]graphs.Id, counter1 int, GlobalValueDataHolder [][]string, GlobalPathDataHolder []RouteInfo) {
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
	var PathHolder RouteInfo
	for l := 0; l < len(path); l++ {
		var EachIDPathHolder []graphs.Id
		var ToStringPathHolder []string
		for k := 0; k < len(path[l]); k++ {
			EachIDPathHolder = path[l]
			ToStringPathHolder = append(ToStringPathHolder, fmt.Sprintf("%v", EachIDPathHolder[k]))
			switch l {
			case 0:
				PathHolder.Route1 = ToStringPathHolder

			case 1:
				PathHolder.Route2 = ToStringPathHolder

			case 2:
				PathHolder.Route3 = ToStringPathHolder
			}
		}
		GlobalPathDataHolder[counter1-1] = PathHolder
	}
}
func Algorithm_Starter(data []Remove_Duplicates.Data) {
	edges := AddEdges(data)
	graph := graphs.NewGraph()
	for i := 0; i < len(edges); i++ {
		err :=graph.AddVertex(edges[i], nil)
		if err !=nil {
			log.Fatal(err)
		}
	}
	// Value from string to float64
	for i := 0; i < len(data); i++ {
		b, err := strconv.ParseFloat(data[i].Value, 64)
		if err != nil {
			log.Fatal(err)
		}

		f, err := strconv.Atoi(data[i].ID)
		if err != nil {
			log.Fatal(err)
		}
		graph.AddEdge(data[i].DepStation, data[i].ArrStation, b, f, err)
		if err!=nil{
			log.Fatal(err)
		}
	}

	var GlobalValueDataHolder [][]string
	var GlobalPathDataHolder []RouteInfo
	var GlobalFromToDataHolder []string
	var counter1 = 0

	GlobalValueDataHolder = make([][]string, len(edges)*(len(edges)-1))
	GlobalPathDataHolder = make([]RouteInfo, len(edges)*(len(edges)-1))
	for i := 0; i < len(edges); i++ {
		for j := 0; j < len(edges); j++ {
			if edges[i] == edges[j] {

			} else {
				counter1++
				dist, path, _ := graph.Yen(edges[i], edges[j], 3)
				fmt.Println("The distance from: "+edges[i]+" to: "+edges[j]+" is ", dist)
				fmt.Println("The path from:"+edges[i]+" to: "+edges[j]+" is: ", path)
				GlobalFromToDataHolder = append(GlobalFromToDataHolder, "from: "+edges[i]+" to: "+edges[j])
				GetPathId(dist, path, counter1, GlobalValueDataHolder, GlobalPathDataHolder)
			}
		}
	}
	var ProcData []ProcessedData
	ProcData = make([]ProcessedData, counter1)
	AddProcdata(ProcData, GlobalValueDataHolder, GlobalPathDataHolder, GlobalFromToDataHolder)
	var NewData InfoArray
	NewData.Array = ProcData
	NewData.DataProcessor(graph)
	get_info_stations(edges, NewData, data)
}
func get_info_stations(edges []string, NewData InfoArray, data []Remove_Duplicates.Data) {
	fmt.Println("You have:", len(edges), "stations")
	fmt.Println(edges)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your current station: ")
	your_station, _ := reader.ReadString('\n')
	your_station = strings.Replace(your_station, "\n", "", -1)

	reader1 := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your destination station: ")
	destination_station, _ := reader1.ReadString('\n')
	destination_station = strings.Replace(destination_station, "\n", "", -1)

	var TextOutPut string
	TextOutPut = "from: " + your_station + " to: " + destination_station
	info_data(TextOutPut, NewData, data)
}
func info_data(TextOutPut string, NewData InfoArray, data []Remove_Duplicates.Data) {
	for _, v := range NewData.Array {
		if TextOutPut == v.FromTo {
			fmt.Println("You have 3 path varieties:")
			fmt.Print("\n")
			//First path

			fmt.Println("First one is: ")
			for _, v := range v.Path.Route1 {
				fmt.Print("->")
				fmt.Print(v)
			}
			fmt.Print("\n")
			fmt.Println("Trains Id's: ")
			for _, v := range v.Path.ID1 {
				fmt.Print("->")
				fmt.Print(v)
			}
			//fmt.Println(v.RouteInfo.ID1)
			fmt.Print("\n")
			fmt.Println("This path will cost for you: ")
			fmt.Println(v.Value[0])
			fmt.Print("\n")
			//Second path

			fmt.Println("Second one is: ")
			for _, v := range v.Path.Route2 {
				fmt.Print("->")
				fmt.Print(v)
			}
			fmt.Println("\nTrains Id's: ")
			for _, v := range v.Path.ID2 {
				fmt.Print("->")
				fmt.Print(v)
			}
			//fmt.Println(v.RouteInfo.ID1)
			fmt.Println("\nThis path will cost for you: ")
			fmt.Println(v.Value[1]) //same
			fmt.Print("\n")
			//Third path

			fmt.Println("Third one is: ")
			for _, v := range v.Path.Route3 {
				fmt.Print("->")
				fmt.Print(v)
			}
			fmt.Println("\nTrains Id's: ")
			for _, v := range v.Path.ID3 {
				fmt.Print("->")
				fmt.Print(v)
			}
			//fmt.Println(v.RouteInfo.ID1)
			fmt.Print("\n")
			fmt.Println("This path will cost for you: ")
			fmt.Println(v.Value[2])
			fmt.Println("Do you want to take more information about trains (y/n) ?")
			reader3 := bufio.NewReader(os.Stdin)
			text3, _ := reader3.ReadString('\n')
			text3 = strings.Replace(text3, "\n", "", -1)
			if text3 == "y" {
				fmt.Println("Information about first route: ")
				for _, v := range v.Path.ID1 {
					for i := 0; i < len(data); i++ {
						a, _ := strconv.Atoi(data[i].ID)

						if v == a {
							fmt.Println(data[i])
						}
					}

				}
				fmt.Println("Information about second route: ")
				for _, v := range v.Path.ID2 {
					for i := 0; i < len(data); i++ {
						a, _ := strconv.Atoi(data[i].ID)

						if v == a {
							fmt.Println(data[i])
						}
					}

				}
				fmt.Println("Information about third route: ")
				for _, k := range v.Path.ID3 {
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
func AddProcdata(ProcData []ProcessedData, GlobalValueDataHolder [][]string, GlobalPathDataHolder []RouteInfo, GlobalFromToDataHolder []string) {
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
		/*for _, v := range ProcData[z].RouteInfo.Route1 {
			fmt.Println(v)
		}*/
	}
	for z := range ProcData {
		ProcData[z].FromTo = GlobalFromToDataHolder[z]
	}
	/*for _, v := range ProcData[z].RouteInfo.Route1 {
		fmt.Println(v)
	}*/
}

func get_handler(handler id_handler, ProcData InfoArray, i int) (id_handler) {
	handler.index1, handler.index2, handler.index3 = -1, -1, -1
	for _, v := range ProcData.Array[i].Path.Route1 {
		handler.id1 = append(handler.id1, v)
		handler.index1++
	}
	for _, v := range ProcData.Array[i].Path.Route2 {
		handler.id2 = append(handler.id2, v)
		handler.index2++
	}
	for _, v := range ProcData.Array[i].Path.Route3 {
		handler.id3 = append(handler.id3, v)
		handler.index3++
	}
	return handler
}
func get_id1(ProcData InfoArray, i int, graph *graphs.Graph, handler id_handler) {
	switch len(handler.id1) {
	case 1:
		fmt.Println("Error") //the same
	case 2:
		ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[handler.id1[handler.index1-1]][handler.id1[handler.index1]].Id)
	case 3:
		ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[handler.id1[handler.index1-2]][handler.id1[handler.index1-1]].Id)
		ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[handler.id1[handler.index1-1]][handler.id1[handler.index1]].Id)
	case 4:
		ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[handler.id1[handler.index1-3]][handler.id1[handler.index1-2]].Id)
		ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[handler.id1[handler.index1-2]][handler.id1[handler.index1-1]].Id)
		ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[handler.id1[handler.index1-1]][handler.id1[handler.index1]].Id)
	case 5:
		ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[handler.id1[handler.index1-4]][handler.id1[handler.index1-3]].Id)
		ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[handler.id1[handler.index1-3]][handler.id1[handler.index1-2]].Id)
		ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[handler.id1[handler.index1-2]][handler.id1[handler.index1-1]].Id)
		ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[handler.id1[handler.index1-1]][handler.id1[handler.index1]].Id)
	}
}
func get_id2(ProcData InfoArray, i int, graph *graphs.Graph, handler id_handler) {
	switch len(handler.id2) {
	case 1:
		fmt.Println("Error")
	case 2:
		ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[handler.id2[handler.index2-1]][handler.id2[handler.index2]].Id)
	case 3:
		ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[handler.id2[handler.index2-2]][handler.id2[handler.index2-1]].Id)
		ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[handler.id2[handler.index2-1]][handler.id2[handler.index2]].Id)
	case 4:
		ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[handler.id2[handler.index2-3]][handler.id2[handler.index2-2]].Id)
		ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[handler.id2[handler.index2-2]][handler.id2[handler.index2-1]].Id)
		ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[handler.id2[handler.index2-1]][handler.id2[handler.index2]].Id)
	case 5:
		ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[handler.id2[handler.index2-4]][handler.id2[handler.index2-3]].Id)
		ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[handler.id2[handler.index2-3]][handler.id2[handler.index2-2]].Id)
		ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[handler.id2[handler.index2-2]][handler.id2[handler.index2-1]].Id)
		ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[handler.id2[handler.index2-1]][handler.id2[handler.index2]].Id)
	}
}
func get_id3(ProcData InfoArray, i int, graph *graphs.Graph, handler id_handler) {
	switch len(handler.id3) {
	case 1:
		fmt.Println("Error")
	case 2:
		ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[handler.id3[handler.index3-1]][handler.id3[handler.index3]].Id)
	case 3:
		ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[handler.id3[handler.index3-2]][handler.id3[handler.index3-1]].Id)
		ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[handler.id3[handler.index3-1]][handler.id3[handler.index3]].Id)
	case 4:
		ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[handler.id3[handler.index3-3]][handler.id3[handler.index3-2]].Id)
		ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[handler.id3[handler.index3-2]][handler.id3[handler.index3-1]].Id)
		ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[handler.id3[handler.index3-1]][handler.id3[handler.index3]].Id)
	case 5:
		ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[handler.id3[handler.index3-4]][handler.id3[handler.index3-3]].Id)
		ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[handler.id3[handler.index3-3]][handler.id3[handler.index3-2]].Id)
		ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[handler.id3[handler.index3-2]][handler.id3[handler.index3-1]].Id)
		ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[handler.id3[handler.index3-1]][handler.id3[handler.index3]].Id)
	}
}
func (ProcData InfoArray) DataProcessor(graph *graphs.Graph) {
	for i := 0; i < len(ProcData.Array); i++ {
		var handler id_handler
		handler = get_handler(handler, ProcData, i)
		get_id1(ProcData, i, graph, handler)
		get_id2(ProcData, i, graph, handler)
		get_id3(ProcData, i, graph, handler)
	}
}
