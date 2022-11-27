package main

import (
	"project/Simple-Go-Project/student/controllers"
	"project/Simple-Go-Project/student/driver"
	"project/Simple-Go-Project/student/models"

	"database/sql"
	"log"
	"net/http"

	//connect to DB istance for elephant sql

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var students []models.Student
var db *sql.DB

func init() {
	gotenv.Load() //loads all the env variable inside the .env file
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db = driver.ConnectDB()
	// log.Println(pgURL)

	/*
		pgURL gives:
		-`dbname`
		-`host`
		-`password`
		-`port`
		-`user`
		pass into sql.open()
	*/

	router := mux.NewRouter()

	controller := controllers.Controller{}

	router.HandleFunc("/students", controller.GetStudents(db)).Methods("GET") //Method-> GET action
	router.HandleFunc("/students/{id}", controller.GetStudent(db)).Methods("GET")
	router.HandleFunc("/students", controller.AddStudent(db)).Methods("GET")
	router.HandleFunc("/students", controller.UpdateStudent(db)).Methods("PUT")
	router.HandleFunc("/students/{id}", controller.RemoveStudent(db)).Methods("DELETE")
	// router.HandleFunc("/students/{id}", getStudent).Methods("GET")
	// router.HandleFunc("/students", addStudent).Methods("POST")
	// router.HandleFunc("/students", updateStudent).Methods("PUT")
	// router.HandleFunc("/students/{id}", removeStudent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
