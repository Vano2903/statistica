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

func NewFileHandler(path, confPath string) FileHandler {
	var f FileHandler
	f.ConfPath = confPath
	f.Path = path
	return f
}

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

func (f FileHandler) UpdateNumOfStudents() error {
	file, err := os.OpenFile(f.Path, os.O_WRONLY|os.O_CREATE, 0744)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(strconv.Itoa(f.NumOfStudents)); err != nil {
		return err
	}
	return nil
}

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

func readNextBytes(file *os.File, recordSize, startFrom int) ([]byte, error) {
	bytes := make([]byte, recordSize)

	_, err := file.ReadAt(bytes, int64(startFrom))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

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
