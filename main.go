package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

const filePath string = "statistica.bin"

type Student struct {
	LastName    [20]rune
	Name        [20]rune
	Phone       [13]rune
	Email       [25]rune
	HasLaptop   bool
	SummerStage bool
}

func (s Student) Encode() ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(s)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (s *Student) Decode(b []byte) error {
	var buff bytes.Buffer
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&s)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var s Student
	s.Name = [20]rune{'a', 'a', 'a', 'a', 'a'}
	s.LastName = [20]rune{'b', 'b', 'b', 'b', 'b'}
	s.Phone = [13]rune{'c', 'c', 'c', 'c', 'c'}
	s.Email = [25]rune{'d', 'd', 'd', 'd', 'd'}
	s.HasLaptop = true
	s.SummerStage = true

	a, _ := s.Encode()
	fmt.Println(a)

	f := NewFileHandler()
	f.Path = filePath

	f.Append(a)

}
