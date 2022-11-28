package studentRepository

import (
	"database/sql"
	"log"
	"project/Simple-Go-Project/student/models"
)

type StudentRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (s StudentRepository) GetStudents(db *sql.DB, student models.Student, students []models.Student) []models.Student {
	rows, err := db.Query("select * from student order by id")
	logFatal(err)

	defer rows.Close()

	for rows.Next() { //map
		err := rows.Scan(&student.ID, &student.Name, &student.Department, &student.DOB) //looking over rows and scanning each record into student; passed as parameter
		logFatal(err)                                                                   //if error

		students = append(students, student)
	}
	return students
}

func (s StudentRepository) GetStudent(db *sql.DB, student models.Student, id int) models.Student {
	rows := db.QueryRow("select * from student where id=$1", id) //$1 => placeholder id inside params
	err := rows.Scan(&student.ID, &student.Name, &student.Department, &student.DOB)
	logFatal(err)

	return student
}

func (s StudentRepository) AddStudent(db *sql.DB, student models.Student) int {
	err := db.QueryRow("insert into student (name, department, dob) values($1, $2, $3) RETURNING id;",
		student.Name, student.Department, student.DOB).Scan(&student.ID) //$-> placeholders
	logFatal(err)

	return student.ID
}

func (s StudentRepository) UpdateStudent(db *sql.DB, student models.Student) int64 {
	result, err := db.Exec("update student set name=$1, department=$2, dob=$3 where id=$4 RETURNING id", &student.Name, &student.Department, &student.DOB, &student.ID)
	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	return rowsUpdated
}

func (s StudentRepository) RemoveStudent(db *sql.DB, id int) int64 {
	result, err := db.Exec("delete from student where id = $1", id) //$1 -> placeholder and the value of placeholder is params["id"]
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	return rowsDeleted
}
