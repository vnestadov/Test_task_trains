package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/vnestadov/Test_task_trains/csv_input"
)
func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter name of your csv file without .csv: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	text = strings.Replace(text, "\n", ".csv", -1)
	csv_input.Input(text)
}
