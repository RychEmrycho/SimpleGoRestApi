package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Student struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Students struct {
	Students []Student `json:"students"`
}

var students Students

func main() {
	jsonFile, err := os.Open("data.json")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully opened data.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &students)

	//for i := 0; i < len(students.Students); i++ {
	//	fmt.Println("ID: " + strconv.Itoa(students.Students[i].ID))
	//	fmt.Println("Name: " + students.Students[i].Name)
	//	fmt.Println("Age: " + strconv.Itoa(students.Students[i].Age))
	//}

	// FOR API
	router := mux.NewRouter()
	router.HandleFunc("/students", GetStudents).Methods("GET")
	router.HandleFunc("/student/{id}", GetStudent).Methods("GET")
	router.HandleFunc("/student/{id}", CreateStudent).Methods("POST")
	router.HandleFunc("/student/{id}", UpdateStudent).Methods("PUT")
	router.HandleFunc("/student/{id}", DeleteStudent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8887", router))
}

func GetStudents(writer http.ResponseWriter, request *http.Request) {
	json.NewEncoder(writer).Encode(students)
}

func GetStudent(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	for _, item := range students.Students {
		if strconv.Itoa(item.ID) == params["id"] {
			json.NewEncoder(writer).Encode(item)
			return
		}
	}
	json.NewEncoder(writer).Encode(&Student{})
}

func CreateStudent(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	var student Student
	_ = json.NewDecoder(request.Body).Decode(&student)
	student.ID, _ = strconv.Atoi(params["id"])
	students.Students = append(students.Students, student)
	json.NewEncoder(writer).Encode(students)
}

func UpdateStudent(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	for index, item := range students.Students {
		if strconv.Itoa(item.ID) == params["id"] {
			_ = json.NewDecoder(request.Body).Decode(&students.Students[index])
			break
		}
	}
	json.NewEncoder(writer).Encode(students)
}

func DeleteStudent(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	for index, item := range students.Students {
		if strconv.Itoa(item.ID) == params["id"] {
			students.Students = append(students.Students[:index], students.Students[index+1:]...)
			break
		}
	}
	json.NewEncoder(writer).Encode(students)
}
