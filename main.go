package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/vnestadov/Test_task_trains/Csv_Input"
)

// Id uniquely identify a vertex.

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter name of your csv file without .csv: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	text = strings.Replace(text, "\n", ".csv", -1)
	Csv_Input.Input(text)
}
