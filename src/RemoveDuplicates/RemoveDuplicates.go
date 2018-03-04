package RemoveDuplicates


import (
	"fmt"

)
type Data struct {
	ID         string
	DepStation string
	ArrStation string
	Value      string
	DepTime    string
	ArrTime    string
}
func RemoveDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}
func RemoveDuplicatesFromData(Data []Data) []Data {
	i := 0
	k := len(Data)
	for _, val := range Data {
		if val.DepStation == Data[i].DepStation && val.ArrStation == Data[i].ArrStation && val.Value == Data[i].Value && val.ArrTime == Data[i].ArrTime && val.DepTime == Data[i].DepTime {
			continue
		}

		Data[i] = val
		i++
	}
	Data = Data[:i]
	j := len(Data)
	fmt.Println("All duplicates deleted from data")
	fmt.Println("Amount of elements before: ", k, "Amount of elements after: ", j)

	return Data
}