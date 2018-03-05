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

	"github.com/vnestadov/Test_task_trains/RemoveDuplicates"
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

//RouteInfo type to save connection between RouteInfo and Id
type RouteInfo struct {
	Route1 []string
	Route2 []string
	Route3 []string
	ID1    []int
	ID2    []int
	ID3    []int
}
func Input(name string) []RemoveDuplicates.Data {
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
func AddEdges (data[]RemoveDuplicates.Data)(edges []string){
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
	edges = RemoveDuplicates.RemoveDuplicatesUnordered(temp)
	return
}
func GetPathId(dist []float64,path [][]graphs.Id,counter1 int,GlobalValueDataHolder [][]string,GlobalPathDataHolder[]RouteInfo){
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
func AlgorithmAccessor(data []RemoveDuplicates.Data) {
	edges := AddEdges(data)
	graph := graphs.NewGraph()
	for i := 0; i < len(edges); i++ {
		graph.AddVertex(edges[i], nil) //todo handle err
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
		graph.AddEdge(data[i].DepStation, data[i].ArrStation, b, f, nil) //todo handle err
	}


	var GlobalValueDataHolder [][]string
	var GlobalPathDataHolder []RouteInfo
	var GlobalFromToDataHolder []string
	var counter1 = 0

	GlobalValueDataHolder = make([][]string, len(edges)*(len(edges)-1))
	GlobalPathDataHolder = make([]RouteInfo, len(edges)*(len(edges)-1))

	//todo so long for, c-style, need refactor it
	for i := 0; i < len(edges); i++ {
		for j := 0; j < len(edges); j++ {
			if edges[i] == edges[j] {

			} else {
				counter1++
				dist, path, _ := graph.Yen(edges[i], edges[j], 3)
				fmt.Println("The distance from: "+edges[i]+" to: "+edges[j]+" is ", dist)
				fmt.Println("The path from:"+edges[i]+" to: "+edges[j]+" is: ", path)
				GlobalFromToDataHolder = append(GlobalFromToDataHolder, "from: "+edges[i]+" to: "+edges[j])
				GetPathId(dist,path,counter1,GlobalValueDataHolder,GlobalPathDataHolder)
			}
		}
	}
	var ProcData []ProcessedData
	ProcData = make([]ProcessedData, counter1)
	AddProcdata(ProcData,GlobalValueDataHolder,GlobalPathDataHolder,GlobalFromToDataHolder)
	var NewData InfoArray
	NewData.Array = ProcData
	NewData.DataProcessor(graph)
	get_info_stations(edges,NewData,data)
}
func get_info_stations(edges[]string,NewData InfoArray,data []RemoveDuplicates.Data) {
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
	info_data(TextOutPut,NewData,data)
}
func info_data(TextOutPut string, NewData InfoArray,data []RemoveDuplicates.Data){
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
//todo to long function
func AddProcdata(ProcData []ProcessedData,GlobalValueDataHolder [][]string,GlobalPathDataHolder[]RouteInfo,GlobalFromToDataHolder[]string){
	for _ = range ProcData { //todo wtf? use just range ProcData
		for h, v := range GlobalValueDataHolder {
			ProcData[h].Value = v

		}
		/*for _, v := range ProcData[z].Value {
			fmt.Println(v)
		}*/
	}
	for _ = range ProcData { //todo the same
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
func (ProcData InfoArray) DataProcessor(graph *graphs.Graph) {
	for i := 0; i < len(ProcData.Array); i++ {
		//var iml RouteInfo
		//todo would be better if you create struct for this and then call function to proccesing this data
		var id1, id2, id3 []string //todo use var(...)
		var index1, index2, index3 int
		index1, index2, index3 = -1, -1, -1 //todo can't understand what is it
		for _, v := range ProcData.Array[i].Path.Route1 {
			id1 = append(id1, v)
			index1++
		}
		for _, v := range ProcData.Array[i].Path.Route2 {
			id2 = append(id2, v)
			index2++
		}
		for _, v := range ProcData.Array[i].Path.Route3 {
			id3 = append(id3, v)
			index3++
		}
		//todo use switch
		if len(id1) == 1 {
			fmt.Println("Error") //the same
		}
		if len(id1) == 2 {

			ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[id1[index1-1]][id1[index1]].Id)
			//fmt.Println("Маршрут 1: ", ProcData.Array[i].RouteInfo.ID1)
		}
		if len(id1) == 3 {
			ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[id1[index1-2]][id1[index1-1]].Id)
			ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[id1[index1-1]][id1[index1]].Id)
			//fmt.Println("Маршрут 1: ", ProcData.Array[i].RouteInfo.ID1)

		}
		if len(id1) == 4 {
			ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[id1[index1-3]][id1[index1-2]].Id)
			ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[id1[index1-2]][id1[index1-1]].Id)
			ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[id1[index1-1]][id1[index1]].Id)
			//fmt.Println("Маршрут 1: ", ProcData.Array[i].RouteInfo.ID1)

		}
		if len(id1) == 5 {
			ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[id1[index1-4]][id1[index1-3]].Id)
			ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[id1[index1-3]][id1[index1-2]].Id)
			ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[id1[index1-2]][id1[index1-1]].Id)
			ProcData.Array[i].Path.ID1 = append(ProcData.Array[i].Path.ID1, graph.Egress[id1[index1-1]][id1[index1]].Id)
			//fmt.Println("Маршрут 1: ", ProcData.Array[i].RouteInfo.ID1)
		}

		//fmt.Println(id1)
		//todo so long create function for this
		if len(id2) == 1 {
			fmt.Println("Error")
		}
		if len(id2) == 2 {

			ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[id2[index2-1]][id2[index2]].Id)
			//fmt.Println("Маршрут 2: ", ProcData.Array[i].RouteInfo.ID2)
		}
		if len(id2) == 3 {
			ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[id2[index2-2]][id2[index2-1]].Id)
			ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[id2[index2-1]][id2[index2]].Id)
			//fmt.Println("Маршрут 2: ", ProcData.Array[i].RouteInfo.ID2)

		}
		if len(id2) == 4 {
			ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[id2[index2-3]][id2[index2-2]].Id)
			ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[id2[index2-2]][id2[index2-1]].Id)
			ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[id2[index2-1]][id2[index2]].Id)
			//fmt.Println("Маршрут 2: ", ProcData.Array[i].RouteInfo.ID2)

		}
		if len(id2) == 5 {
			ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[id2[index2-4]][id2[index2-3]].Id)
			ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[id2[index2-3]][id2[index2-2]].Id)
			ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[id2[index2-2]][id2[index2-1]].Id)
			ProcData.Array[i].Path.ID2 = append(ProcData.Array[i].Path.ID2, graph.Egress[id2[index2-1]][id2[index2]].Id)
			//fmt.Println("Маршрут 2: ", ProcData.Array[i].RouteInfo.ID2)
		}
		//fmt.Println(id2)
		//todo the same
		if len(id3) == 1 {
			fmt.Println("Error")
		}
		if len(id3) == 2 {

			ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[id3[index3-1]][id3[index3]].Id)
			//fmt.Println("Маршрут 3: ", ProcData.Array[i].RouteInfo.ID3)
		}
		if len(id3) == 3 {
			ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[id3[index3-2]][id3[index3-1]].Id)
			ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[id3[index3-1]][id3[index3]].Id)
			//fmt.Println("Маршрут 3: ", ProcData.Array[i].RouteInfo.ID3)

		}
		if len(id3) == 4 {
			ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[id3[index3-3]][id3[index3-2]].Id)
			ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[id3[index3-2]][id3[index3-1]].Id)
			ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[id3[index3-1]][id3[index3]].Id)
			//fmt.Println("Маршрут 3: ", ProcData.Array[i].RouteInfo.ID3)

		}
		if len(id3) == 5 {
			ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[id3[index3-4]][id3[index3-3]].Id)
			ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[id3[index3-3]][id3[index3-2]].Id)
			ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[id3[index3-2]][id3[index3-1]].Id)
			ProcData.Array[i].Path.ID3 = append(ProcData.Array[i].Path.ID3, graph.Egress[id3[index3-1]][id3[index3]].Id)
			//fmt.Println("Маршрут 3: ", ProcData.Array[i].RouteInfo.ID3)
		}
		//fmt.Println(id3)

	}

}
