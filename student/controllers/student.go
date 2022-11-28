package controllers

import (
	"project/Simple-Go-Project/student/models"
	studentRepository "project/Simple-Go-Project/student/repository"
	"strconv"

	"database/sql"
	"log"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Controller struct{}

var students []models.Student

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c Controller) GetStudents(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { //returning a callback function
		var student models.Student
		students = []models.Student{}

		studentRepo := studentRepository.StudentRepository{} //StudentRepository struct
		students = studentRepo.GetStudents(db, student, students)
		json.NewEncoder(w).Encode(students)
	}
}

func (c Controller) GetStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student models.Student
		params := mux.Vars(r)

		students = []models.Student{}
		studentRepo := studentRepository.StudentRepository{}

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		student = studentRepo.GetStudent(db, student, id)

		json.NewEncoder(w).Encode(student)
	}
}

func (c Controller) AddStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student models.Student
		var studentID int

		//Decode -> values inside the request body are mapped to the fields of student object
		json.NewDecoder(r.Body).Decode(&student)

		studentRepo := studentRepository.StudentRepository{}
		studentID = studentRepo.AddStudent(db, student)

		json.NewEncoder(w).Encode(studentID)
	}
}

func (c Controller) UpdateStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student models.Student
		json.NewDecoder(r.Body).Decode(&student) //decode request body and mapping

		studentRepo := studentRepository.StudentRepository{}
		rowsUpdated := studentRepo.AddStudent(db, student)

		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) RemoveStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r) //returns key value pairs-> id passed in the url

		studentRepo := studentRepository.StudentRepository{}

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		rowsDeleted := studentRepo.RemoveStudent(db, id)
		json.NewEncoder(w).Encode(rowsDeleted)
	}
}
