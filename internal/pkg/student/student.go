package student

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Student struct {
	LastName    [30]byte
	Name        [30]byte
	Phone       [13]byte
	Email       [50]byte
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
	//create a decode
	dec := gob.NewDecoder(&buff)
	//decoding the slice of byte
	err := dec.Decode(&s)
	if err != nil {
		return err
	}
	return nil
}

//print a student in a formatted way
func (s Student) Print() {
	//range is like a forEach and the variable "a" is a byte so we convert it to string and print it

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
