package main

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"strconv"
	"unsafe"
)

type FileHandler struct {
	Path          string
	ConfPath      string
	NumOfStudents int
}

//"constructor"
func NewFileHandler(path, confPath string) FileHandler {
	var f FileHandler
	f.ConfPath = confPath
	f.Path = path
	return f
}

//get the number of students from the configuration file
func (f *FileHandler) GetNumOfStudents() error {
	content, err := ioutil.ReadFile(f.ConfPath)
	if err != nil {
		return err
	}

	f.NumOfStudents, err = strconv.Atoi(string(content))
	if err != nil {
		return err
	}
	return nil
}

//write the number of students in the configuration file
func (f FileHandler) UpdateNumOfStudents() error {
	file, err := os.OpenFile(f.ConfPath, os.O_WRONLY|os.O_CREATE, 0744)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(strconv.Itoa(f.NumOfStudents)); err != nil {
		return err
	}
	return nil
}

//append a student in the file, also update the configuration file
func (f *FileHandler) Append(b []byte) error {
	file, err := os.OpenFile(f.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0744)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(b); err != nil {
		return err
	}
	f.NumOfStudents++
	f.UpdateNumOfStudents()
	return nil
}

//read from a file a record starting from an offset (int)
func readNextBytes(file *os.File, recordSize, startFrom int) ([]byte, error) {
	bytes := make([]byte, recordSize)

	_, err := file.ReadAt(bytes, int64(startFrom))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

//return a slice with all the students stored in the file
func (f FileHandler) GetAllStudents() ([]Student, error) {
	var students []Student
	var temp Student
	size := int(unsafe.Sizeof(temp))
	file, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = f.UpdateNumOfStudents()
	if err != nil {
		return nil, err
	}
	for i := 0; i < f.NumOfStudents; i++ {
		var s Student
		data, err := readNextBytes(file, size, size*i)
		if err != nil {
			return nil, err
		}
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &s)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}
