package fileHandler

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
	"unsafe"

	"github.com/Vano2903/statistica/internal/pkg/student"
	"github.com/Vano2903/statistica/internal/pkg/utils/input"
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
	//Stats can return a error saying if the file doesn't exist, check if the error
	//returned is because the file doesn't exist, if so create and initialize it
	if _, err := os.Stat(f.Path); os.IsNotExist(err) {
		f.UpdateNumOfStudents()
		return nil
	}

	//open the file
	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	//this function will run before the function returns (will close the file stream)
	defer file.Close()

	//get the size of a uint32 and read this size from the start of the file
	var temp uint32
	content, err := readNextBytes(file, int(unsafe.Sizeof(temp)), 0)
	if err != nil {
		return err
	}

	//convert the content of the file in uint32
	f.NumOfStudents = binary.BigEndian.Uint32(content)
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

	//write on the file the slice of byte given in argument
	if _, err := file.Write(b); err != nil {
		return err
	}

	//close the file stream, not using defer because UpdateNumOfStudents() open the same file
	file.Close()

	//increase the number of students
	f.NumOfStudents++
	//update the number of students in the file
	f.UpdateNumOfStudents()
	return nil
}

//read from a file a record starting from an offset (int)
func readNextBytes(file *os.File, recordSize, startFrom int) ([]byte, error) {
	//create a slice of byte with the record's size
	bytes := make([]byte, recordSize)

	//read for the length of the record starting after an offset ammont of bytes
	_, err := file.ReadAt(bytes, int64(startFrom))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

//given the number of a record will return the student inside it
func (f FileHandler) GetStudent(recordNum int) (student.Student, error) {
	var s student.Student
	//size of student stuct
	size := int(unsafe.Sizeof(s))
	//open file
	file, err := os.Open(f.Path)
	//the stream will close right before the function will return
	defer file.Close()
	if err != nil {
		return student.Student{}, err
	}

	var temp int32

	//the size of the record is the student*the record asked + size of uint32
	data, err := readNextBytes(file, size, (size*recordNum)+int(unsafe.Sizeof(temp)))
	if err != nil {
		return student.Student{}, err
	}

	//read the []byte returned by the readNextBytes() and fill the s (student) variable
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
	var i int
	for {
		//get a student given the for loop counter
		s, err := f.GetStudent(i)
		//if a error will return, it means its EOF
		if err != nil {
			return students, nil
		}
		//append to slice of student
		students = append(students, s)
		i++
	}
}

//given a phone number it will return the student with that phone (it's a primary key)
func (f FileHandler) SearchByPhone(phone string) (student.Student, error) {
	var s student.Student
	var tempUint uint32

	//size of the array of byte that contains the phone
	PhoneSize := int(unsafe.Sizeof(s.Phone))
	//ammount of bytes to get to the first phone record in the file
	first := int(unsafe.Sizeof(s.LastName)) + int(unsafe.Sizeof(s.Name)) + int(unsafe.Sizeof(tempUint))
	//ammount of bytes from a phone record to another
	others := first + int(unsafe.Sizeof(student.Student{}))

	//open the file stream
	file, err := os.Open(f.Path)
	if err != nil {
		return student.Student{}, err
	}
	//the stream will close right before the function will return
	defer file.Close()

	var i int
	for {
		var phoneByte []byte
		if i == 0 {
			//first phone number
			phoneByte, err = readNextBytes(file, PhoneSize, first)
		} else {
			//others phone numbers
			phoneByte, err = readNextBytes(file, PhoneSize, others*i)
		}
		//remove all the unused bytes of the array
		phoneByte = bytes.Trim(phoneByte, string('\u0000'))

		//check if its EOF
		if err != nil {
			return student.Student{}, errors.New("phone number was not found")
		}

		//convert the phone record to string and check if its the same to the phone given as argument
		if phone == string(phoneByte) {
			//if it's the same get the student inside the record we are checking
			s, err := f.GetStudent(i)
			if err != nil {
				return student.Student{}, err
			}

			return s, nil
		}
		i++
	}
}

func (f FileHandler) SearchByNames(lastname, name string) (student.Student, error) {
	archiveIndex := f.Hash_3(name, lastname)
	newStudent, err := f.GetStudent(archiveIndex)
	if err != nil {
		return student.Student{}, err
	}
	//remove all unused bytes, convert the parameter records to string and check if its the same to the lastname and name given as argument
	if string(bytes.Trim(newStudent.LastName[:], string('\u0000'))) == lastname && string(bytes.Trim(newStudent.Name[:], string('\u0000'))) == name {
		return newStudent, nil
	}

	cf := NewFileHandler("collisions.bin")
	if err := cf.GetNumOfStudents(); err != nil {
		return student.Student{}, err
	}

	//open the file stream
	file, err := os.Open(cf.Path)
	//the stream will close right before the function will return
	defer file.Close()
	if err != nil {
		return student.Student{}, err
	}

	for {
		var collision struct {
			collidedIndex uint16
			newIndex      uint16
		}
		var tempUint uint32

		recordBytes, err := readNextBytes(file, int(unsafe.Sizeof(collision)), int(unsafe.Sizeof(tempUint)))
		if err != nil {
			return student.Student{}, err
		}
		collision.collidedIndex = binary.BigEndian.Uint16(recordBytes[:len(recordBytes)/2])
		collision.newIndex = binary.BigEndian.Uint16(recordBytes[:len(recordBytes)/2])

		if collision.collidedIndex != uint16(archiveIndex) {
			continue
		}
		checkStudent, err := f.GetStudent(int(collision.newIndex))
		if err != nil {
			return student.Student{}, err
		}
		//remove all unused bytes, convert the parameter records to string and check if its the same to the lastname and name given as argument
		if string(bytes.Trim(checkStudent.LastName[:], string('\u0000'))) == lastname && string(bytes.Trim(checkStudent.Name[:], string('\u0000'))) == name {
			return checkStudent, nil
		}
	}
}

func (f FileHandler) AddStudent() {
	var s student.Student

	for {
		fmt.Printf("nome dello studente [%d]: ", len(s.Name))
		name := input.String()
		if len(name) <= len(s.Name) {
			copy(s.Name[:], name)
			break
		} else {
			fmt.Printf("\nnome troppo lungo, massimo %d caratteri\n", len(s.Name))
		}
	}

	for {
		fmt.Printf("cognome dello studente [%d]: ", len(s.LastName))
		lastName := input.String()
		if len(lastName) <= len(s.LastName) {
			copy(s.LastName[:], lastName)
			break
		} else {
			fmt.Printf("\ncognome troppo lungo, massimo %d caratteri\n", len(s.LastName))
		}
	}

	for {
		fmt.Printf("telefono dello studente [%d]: ", len(s.Phone))
		phone := input.String()
		if len(phone) <= len(s.Phone) {
			copy(s.Phone[:], phone)
			break
		} else {
			fmt.Printf("\ntelefono troppo lungo, massimo %d caratteri\n", len(s.Phone))
		}
	}

	for {
		fmt.Printf("email dello studente [%d]: ", len(s.Email))
		email := input.String()
		if len(email) <= len(s.Email) {
			copy(s.Email[:], email)
			break
		} else {
			fmt.Printf("\nemail troppo lungo, massimo %d caratteri\n", len(s.Email))
		}
	}

	for {
		fmt.Print("ha un laptop [t/f]: ")
		laptop := input.String()
		if laptop == "t" || laptop == "f" {
			s.HasLaptop = laptop[0]
			break
		} else {
			fmt.Println("\nil valore inserito non é accettato, usa t (true) f (false)")
		}
	}

	for {
		fmt.Print("ha fatto lo stage [t/f]: ")
		stage := input.String()
		if stage == "t" || stage == "f" {
			s.SummerStage = stage[0]
			break
		} else {
			fmt.Println("\nil valore inserito non é accettato, usa t (true) f (false)")
		}
	}

	//f.Append(s.Encode())
	f.WriteStudent(s.Encode(), f.Hash_3(string(s.Name[:]), string(s.LastName[:])))
}

func (f *FileHandler) WriteStudent (b []byte, index int) (error){
	//open file in read and write mode
	file, err := os.OpenFile(f.Path, os.O_RDWR, 0744)
	//the stream will close right before the function will return
	defer file.Close()
	if err != nil {
		return err
	}

	//Check if the index position is a collision
	var tempStudent student.Student
	var tempUint uint32
	recordBytes, err := readNextBytes(file, int(unsafe.Sizeof(tempStudent)), int(unsafe.Sizeof(tempUint))+int(unsafe.Sizeof(tempStudent))*index)
	var isACollision bool
	for _, byte := range recordBytes {
		if byte != 0 {
			isACollision = true
			break
		}
	}

	//write on the file at the index position
	if isACollision == false {
		_, err := file.WriteAt(b, int64(int(unsafe.Sizeof(tempUint))+int(unsafe.Sizeof(tempStudent))*index))
		f.NumOfStudents++
		f.UpdateNumOfStudents()
		if err != nil {
			return err
		}
	}

	//find the first free record on the file
	var freeRecord bool
	for s := 0; s < int(f.NumOfStudents); s++ {
		checkRecord, err := readNextBytes(file, int(unsafe.Sizeof(tempStudent)), int(unsafe.Sizeof(tempUint))+int(unsafe.Sizeof(tempStudent))*s)
		if err != nil {
			return err
		}

		freeRecord = true
		for _, byte := range checkRecord {
			if byte != 0 {
				freeRecord = false
				break
			}
		}

		if freeRecord == false {
			continue
		}

		//Write the student record on the file at new index
		_, err = file.WriteAt(b, int64(int(unsafe.Sizeof(tempUint))+int(unsafe.Sizeof(tempStudent))*s))
		f.NumOfStudents++
		f.UpdateNumOfStudents()
		if err != nil {
			return err
		}

		//Generate the handler for the collisions file
		cf := NewFileHandler("collisions.bin")
		if err := cf.GetNumOfStudents(); err != nil {
			return err
		}
		//open the file stream
		collisionsFile, err := os.OpenFile(cf.Path, os.O_WRONLY|os.O_APPEND, 0744)
		//the stream will close right before the function will return
		defer collisionsFile.Close()
		if err != nil {
			return err
		}

		collidedIndex := []byte(strconv.Itoa(index))
		newIndex := []byte(strconv.Itoa(s))
		_, err = collisionsFile.Write(collidedIndex)
		if err != nil {
			return err
		}
		_, err = collisionsFile.Write(newIndex)
		if err != nil {
			return err
		}
		cf.NumOfStudents++
		cf.UpdateNumOfStudents()
		return nil
	}

	return fmt.Errorf("archivio pieno: impossibile salvare un nuovo dato")
}