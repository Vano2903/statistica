package fileHandler

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"unsafe"

	"github.com/Vano2903/statistica-go/internal/pkg/student"
)

type FileHandler struct {
	Path          string
	NumOfStudents uint32
}

//"constructor"
func NewFileHandler(path string) FileHandler {
	var f FileHandler
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

func (f FileHandler) GetStudent(recordNum int) (student.Student, error) {
	var s student.Student
	size := int(unsafe.Sizeof(s))
	//open file
	file, err := os.Open(f.Path)
	if err != nil {
		return student.Student{}, err
	}
	//the stream will close right before the function will return
	defer file.Close()

	var temp int32

	data, err := readNextBytes(file, size, (size*recordNum)+int(unsafe.Sizeof(temp)))
	if err != nil {
		return student.Student{}, err
	}
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		return student.Student{}, err
	}
	return s, nil
}

//return a slice with all the students stored in the file
func (f FileHandler) GetAllStudents() ([]student.Student, error) {
	var students []student.Student
	//f.GetNumOfStudents()
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

func (f FileHandler) SearchByPhone(phone string) (student.Student, error) {
	type temp struct {
		LastName [20]byte
		Name     [20]byte
	}
	type temp2 struct {
		Email       [25]byte
		HasLaptop   byte
		SummerStage byte
	}
	var s student.Student
	var tempUint uint32
	var i int
	PhoneSize := int(unsafe.Sizeof([20]byte{}))
	file, err := os.Open(f.Path)
	if err != nil {
		return student.Student{}, err
	}
	//the stream will close right before the function will return
	defer file.Close()
	for {
		var phoneByte []byte
		if i == 0 {
			phoneByte, err = readNextBytes(file, PhoneSize, int(unsafe.Sizeof(temp{})))
		} else {
			phoneByte, err = readNextBytes(file, PhoneSize, int(unsafe.Sizeof(temp2{}))+int(unsafe.Sizeof(temp{}))+int(unsafe.Sizeof(tempUint)))
		}
		if err != nil {
			return student.Student{}, errors.New("phone number was not found")
		}
		if phone == string(phoneByte) {
			return s, nil
		}
		i++
	}
}