package main

import (
	"fmt"
	"log"

	"github.com/Vano2903/statistica/internal/pkg/fileHandler"
	"github.com/Vano2903/statistica/internal/pkg/utils/clear"
	"github.com/Vano2903/statistica/internal/pkg/utils/input"
)

const filePath string = "statistica.bin"

var f fileHandler.FileHandler

//initialize the file handler and get the number of students
func init() {
	f = fileHandler.NewFileHandler(filePath)
	if err := f.GetNumOfStudents(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	for {
		clear.Clear()
		fmt.Println("STATISTICA - VANONCINI && MORANDI")
		fmt.Println("1] aggiungi uno studente")
		fmt.Println("2] visualizza tutti gli studenti")
		fmt.Println("3] cerca uno studente per telefono")
		fmt.Println("4] cerca uno studente per nome e cognome")
		fmt.Println("5] esci dal programma")
		choice := input.String()
		switch choice {
		case "1":
			f.AddStudent()
		case "2":
			stu, err := f.GetAllStudents()
			if err != nil {
				log.Fatal(err)
			}
			for _, a := range stu {
				a.Print()
				fmt.Println("")
			}
		case "3":
			fmt.Print("numero di telefon: ")
			phone := input.String()
			s, err := f.SearchByPhone(phone)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				s.Print()
			}
		case "4":
			//! cerca per nome e cognome
		case "5":
			fmt.Println("bye")
			return
		}
	}
}
