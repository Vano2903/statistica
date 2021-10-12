package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"unsafe"
)

type FileHandler struct {
	Path          string
	ConfPath      string
	NumOfStudents uint32
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
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		f.UpdateNumOfStudents()
		return nil
	}

	var temp uint32
	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	content, err := readNextBytes(file, int(unsafe.Sizeof(temp)), 0)

	//convert the content of the file in int
	f.NumOfStudents = binary.BigEndian.Uint32(content)
	if err != nil {
		return err
	}
	return nil
}

//write the number of students in the configuration file
func (f FileHandler) UpdateNumOfStudents() error {
	//open the file in write only and create if not exist
	file, err := os.OpenFile(f.Path, os.O_WRONLY|os.O_CREATE, 0744)
	if err != nil {
		return err
	}
	//close the stream before the function returns
	defer file.Close()

	//convert the number of students in slice of byte
	numStudentsBinary := make([]byte, 4)
	binary.LittleEndian.PutUint32(numStudentsBinary, f.NumOfStudents)

	//write the slice of byte in the file
	if _, err := file.Write(numStudentsBinary); err != nil {
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

func (f FileHandler) GetStudent(recordNum int) (Student, error) {
	var s Student
	size := int(unsafe.Sizeof(s))
	//open file
	file, err := os.Open(f.Path)
	if err != nil {
		return Student{}, err
	}
	//the stream will close right before the function will return
	defer file.Close()

	var temp int32

	data, err := readNextBytes(file, size, (size*recordNum)+int(unsafe.Sizeof(temp)))
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

//return a slice with all the students stored in the file
func (f FileHandler) GetAllStudents() ([]Student, error) {
	var students []Student
	f.GetNumOfStudents()
	//needed just for the sizeof
	var i int
	for true {
		s, err := f.GetStudent(i)
		if err != nil {
			return students, nil
		}
		students = append(students, s)
		i++
	}
	return students, nil
}

func (f FileHandler) SearchByPhone(phone string) ([]Student, error) {
	return nil, nil
}
