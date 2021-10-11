package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"unsafe"
)

type FileHandler struct {
	Path string
}

func NewFileHandler() FileHandler {
	var f FileHandler
	return f
}

func (f FileHandler) Append(b []byte) error {
	file, err := os.OpenFile(f.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0744)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(b); err != nil {
		return err
	}
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

func (f FileHandler) GetAllStudents() (Student, error) {
	var s Student
	size := int(unsafe.Sizeof(s))
	file, err := os.Open(f.Path)
	if err != nil {
		return Student{}, err
	}
	defer file.Close()

	data, err := readNextBytes(file, size, size*1)
	if err != nil {
		return Student{}, err
	}
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		return Student{}, err
	}

	return s, nil
}
