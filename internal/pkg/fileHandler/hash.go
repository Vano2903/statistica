package fileHandler

/*
type Path string
type LogicRecord uint32
type FileHandler struct {
	path Path
	size LogicRecord
}
*/

func (file FileHandler) Hash_1(firstname, lastname string) int {
	hashableString := firstname + lastname

	var sum int
	for index, character := range hashableString {
		sum += int(character * rune(index+1))
	}

	return (sum / len(hashableString)) % int(file.NumOfStudents)
}

func (file FileHandler) Hash_2(firstname, lastname string) int {
	hashableString := firstname + lastname

	var sum int
	for index, character := range hashableString {
		sum += int(character*rune(index+1)) / len(hashableString)
	}

	return sum % int(file.NumOfStudents)
}

func (file FileHandler) Hash_3(firstname, lastname string) int {
	hashableString := firstname + lastname

	var sum int
	for index, character := range hashableString {
		sum += int(character * rune(index+1))
	}

	return sum % int(file.NumOfStudents)
}
