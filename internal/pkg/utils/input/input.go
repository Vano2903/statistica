package input

import (
	"bufio"
	"os"
	"strconv"
)

/* -- Int type -- */

//Input a string and convert it to an integer
func Int() int {
	input := getString()
	inputInt, _ := strconv.Atoi(input)
	return inputInt
}

//Asks the integer value until it complies with the limits set
func ControlledInt(min, max int) int {
	for {
		input := getString()
		inputInt, _ := strconv.Atoi(input)
		if inputInt >= min && inputInt <= max {
			return inputInt
		}
	}
}

//Compare the integer value taken as input with the value 0 to return the result of the comparison
func CompareInt() int {
	inputInt := Int()
	switch {
	case inputInt < 0:
		return -1
	case inputInt == 0:
		return 0
	case inputInt > 0:
		return 1
	}
	return inputInt
}

/* -- Float type -- */

//Input a string and convert it to an float64
func Float() float64 {
	input := getString()
	inputFloat, _ := strconv.ParseFloat(input, 64)
	return inputFloat
}

/* -- String type -- */

//Ask for a string as input
func String() string {
	input := getString()
	return input
}

/* -- Boolean type -- */

//Input a string and convert it to Boolean
func Bool() bool {
	input := getString()
	if input == "false" || input == "0" {
		return false
	}
	return true
}

/* -- Kirito "get" function -- */

//Input a string
func getString() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	string := scanner.Text()
	return string
}
