package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

const filePath string = "statistica.bin"

var f FileHandler

type Student struct {
	LastName    [20]byte
	Name        [20]byte
	Phone       [13]byte
	Email       [25]byte
	HasLaptop   byte
	SummerStage byte
}

//convert the student struct to slice of byte
func (s Student) Encode() []byte {
	var complete []byte
	for _, b := range s.LastName {
		complete = append(complete, b)
	}
	for _, b := range s.Name {
		complete = append(complete, b)
	}
	for _, b := range s.Phone {
		complete = append(complete, b)
	}
	for _, b := range s.Email {
		complete = append(complete, b)
	}
	complete = append(complete, s.HasLaptop)
	complete = append(complete, s.SummerStage)
	return complete
}

//given a slice of byte it will try to convert it to the student struct
func (s *Student) Decode(b []byte) error {
	var buff bytes.Buffer
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&s)
	if err != nil {
		return err
	}
	return nil
}

//print a student in a formatted way
func (s Student) Print() {
	fmt.Print("nome: ")
	for _, a := range s.Name {
		fmt.Print(string(a))
	}
	fmt.Println("")

	fmt.Print("cognome: ")
	for _, a := range s.LastName {
		fmt.Print(string(a))
	}
	fmt.Println("")

	fmt.Print("telefono: ")
	for _, a := range s.Phone {
		fmt.Print(string(a))
	}
	fmt.Println("")

	fmt.Print("email: ")
	for _, a := range s.Email {
		fmt.Print(string(a))
	}
	fmt.Println("")

	fmt.Print("ha un laptop: ")
	if string(s.HasLaptop) == "t" {
		fmt.Print("si")
	} else {
		fmt.Print("no")
	}
	fmt.Println("")

	fmt.Print("ha fatto lo stage: ")
	if string(s.SummerStage) == "t" {
		fmt.Print("si")
	} else {
		fmt.Print("no")
	}
	fmt.Println("")
}

//initialize the file handler and get the number of students
func init() {
	f = NewFileHandler(filePath)
	if err := f.GetNumOfStudents(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// var s Student
	// s.Name = [20]byte{'a', 'a', 'a', 'a', 'a'}
	// s.LastName = [20]byte{'b', 'b', 'b', 'b', 'b'}
	// s.Phone = [13]byte{'c', 'c', 'c', 'c', 'c'}
	// s.Email = [25]byte{'d', 'd', 'd', 'd', 'd'}
	// s.HasLaptop = byte('t')
	// s.SummerStage = byte('f')

	// a := s.Encode()

	// s.Name = [20]byte{'A', 'A', 'A', 'A', 'A'}
	// s.LastName = [20]byte{'B', 'B', 'B', 'B', 'B'}
	// s.Phone = [13]byte{'C', 'C', 'C', 'C', 'C'}
	// s.Email = [25]byte{'D', 'D', 'D', 'D', 'D'}
	// s.HasLaptop = byte('f')
	// s.SummerStage = byte('t')

	// f := NewFileHandler(filePath, confFilePath)
	// f.Path = filePath

	// f.Append(a)

	// b := s.Encode()
	// f.Append(b)

	stu, err := f.GetAllStudents()
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range stu {
		a.Print()
		fmt.Println("")
	}

}
