package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"CsvInput"
)



// Id uniquely identify a vertex.



func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter name of your csv file without .csv: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	text = strings.Replace(text, "\n", ".csv", -1)

	CsvInput.CsvInputFunc(text)
}


