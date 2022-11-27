package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"project/Simple-Go-Project/student/models"

	//connect to DB istance for elephant sql
	"github.com/lib/pq"
	"github.com/subosito/gotenv"

	"github.com/gorilla/mux"
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
	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	logFatal(err)

	db, err = sql.Open("postgres", pgURL)

	logFatal(err)

	err = db.Ping()
	logFatal(err) //if ping() gives error

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

	router.HandleFunc("/students", getStudents).Methods("GET") //Method-> GET action
	router.HandleFunc("/students/{id}", getStudent).Methods("GET")
	router.HandleFunc("/students", addStudent).Methods("POST")
	router.HandleFunc("/students", updateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", removeStudent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	students = []models.Student{}
	rows, err := db.Query("select * from student order by id")
	logFatal(err)

	defer rows.Close()

	for rows.Next() { //map
		err := rows.Scan(&student.ID, &student.Name, &student.Department, &student.DOB)
		logFatal(err) //if error

		students = append(students, student)
	}

	json.NewEncoder(w).Encode(students)

}

func getStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	params := mux.Vars(r)

	rows := db.QueryRow("select * from student where id=$1", params["id"]) //$1 => placeholder id inside params
	err := rows.Scan(&student.ID, &student.Name, &student.Department, &student.DOB)
	logFatal(err)

	json.NewEncoder(w).Encode(student)
}

func addStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	var studentID int

	//Decode -> values inside the request body are mapped to the fields of student object
	json.NewDecoder(r.Body).Decode(&student)
	err := db.QueryRow("insert into student (name, department, dob) values($1, $2, $3) RETURNING id;",
		student.Name, student.Department, student.DOB).Scan(&studentID) //$-> placeholders
	logFatal(err)
	json.NewEncoder(w).Encode(studentID)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	json.NewDecoder(r.Body).Decode(&student) //decode request body and mapping

	result, err := db.Exec("update student set name=$1, department=$2, dob=$3 where id=$4 RETURNING id", &student.Name, &student.Department, &student.DOB, &student.ID)
	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)
}

func removeStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)                                                     //returns key value pairs-> id passed in the url
	result, err := db.Exec("delete from student where id = $1", params["id"]) //$1 -> placeholder and the value of placeholder is params["id"]
	logFatal(err)

	roewsDeleted, err := result.RowsAffected()
	json.NewEncoder(w).Encode(roewsDeleted)
}
