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
	//create the file if not exist and initialize it with 0
	if _, err := os.Stat(f.ConfPath); os.IsNotExist(err) {
		f.UpdateNumOfStudents()
		return nil
	}
	//read the config file
	content, err := ioutil.ReadFile(f.ConfPath)
	if err != nil {
		return err
	}

	//convert the content of the file in int
	f.NumOfStudents, err = strconv.Atoi(string(content))
	if err != nil {
		return err
	}
	return nil
}

//write the number of students in the configuration file
func (f FileHandler) UpdateNumOfStudents() error {
	//open the file in write only and create if not exist
	file, err := os.OpenFile(f.ConfPath, os.O_WRONLY|os.O_CREATE, 0744)
	if err != nil {
		return err
	}
	//close the stream before the function returns
	defer file.Close()

	//write to the file the string convetion of the number of students and check the error
	if _, err := file.WriteString(strconv.Itoa(f.NumOfStudents)); err != nil {
		return err
	}
	return nil
}

//append a student in the file, also update the configuration file
func (f *FileHandler) Append(b []byte) error {
	//open file in write only, create if not exist and append mode
	file, err := os.OpenFile(f.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0744)
	if err != nil {
		return err
	}
	//close the stream before the function returns
	defer file.Close()

	//write on the file the slice of byte
	if _, err := file.Write(b); err != nil {
		return err
	}
	//increase the number of students
	f.NumOfStudents++
	//update in the configuration file the number of students
	f.UpdateNumOfStudents()
	return nil
}

//read from a file a record starting from an offset (int)
func readNextBytes(file *os.File, recordSize, startFrom int) ([]byte, error) {
	//create a slice of byte with the record's size
	bytes := make([]byte, recordSize)

	//read for the length of a record starting after an offset ammont of bytes
	_, err := file.ReadAt(bytes, int64(startFrom))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

//return a slice with all the students stored in the file
func (f FileHandler) GetAllStudents() ([]Student, error) {
	var students []Student
	//needed just for the sizeof
	var temp Student
	size := int(unsafe.Sizeof(temp))

	//open file
	file, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	//the stream will close right before the function will return
	defer file.Close()

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
