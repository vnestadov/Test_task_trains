package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vnestadov/Test_task_trains/CsvInput"
)

// Id uniquely identify a vertex.

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter name of your csv file without .csv: ") //todo read file directly from program or create config for this and read him from program args
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	text = strings.Replace(text, "\n", ".csv", -1)

	CsvInput.CsvInputFunc(text)
}
