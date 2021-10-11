package main

import "os"

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

func (f FileHandler) GetAllStudents() ([]Student, error) {
	return nil, nil
}
