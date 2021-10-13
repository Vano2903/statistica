package hash

import "github.com/Vano2903/statistica-go/internal/pkg/fileHandler"

/*
type Path string
type LogicRecord uint32
type FileHandler struct {
	path Path
	size LogicRecord
}
*/

func (file fileHandler.FileHandler) Hash_1(firstname, lastname string) int {
	hashableString := firstname + lastname

	var sum int
	for index, character := range hashableString {
		sum += int(character * rune(index+1))
	}

	return (sum / len(hashableString)) % int(file.NumOfStudents)
}

func (file fileHandler.FileHandler) Hash_2(firstname, lastname string) int {
	hashableString := firstname + lastname

	var sum int
	for index, character := range hashableString {
		sum += int(character*rune(index+1)) / len(hashableString)
	}

	return sum % int(file.NumOfStudents)
}

func (file fileHandler.FileHandler) Hash_3(firstname, lastname string) int {
	hashableString := firstname + lastname

	var sum int
	for index, character := range hashableString {
		sum += int(character * rune(index+1))
	}

	return sum % int(file.NumOfStudents)
}
