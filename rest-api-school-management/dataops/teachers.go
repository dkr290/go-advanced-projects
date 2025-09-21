package dataops

import (
	"database/sql"
	"fmt"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

type DatabaseInf interface {
	InsertTeachers(*models.Teacher) (int64, error)
	GetTeacherByID(int) (models.Teacher, error)
	GetAllTeachers(string, string) (*sql.Rows, error)
}

type Teachers struct {
	db *sql.DB
}

func NewTeachersDB(db *sql.DB) *Teachers {
	return &Teachers{
		db: db,
	}
}

func (t *Teachers) InsertTeachers(tm *models.Teacher) (int64, error) {
	stmt, err := t.db.Prepare(`INSERT INTO teachers
		            (first_name,last_name,email,class,subject)
                VALUES(?,?,?,?,?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	sqlResp, err := stmt.Exec(
		tm.FirstName,
		tm.LastName,
		tm.Email,
		tm.Class,
		tm.Subject,
	)
	if err != nil {
		return 0, err
	}
	lastID, err := sqlResp.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

func (t *Teachers) GetTeacherByID(id int) (models.Teacher, error) {
	var teacher models.Teacher
	err := t.db.QueryRow("SELECT id, first_name, last_name ,email, class, subject FROM teachers WHERE id = ?", id).
		Scan(
			&teacher.ID,
			&teacher.FirstName,
			&teacher.LastName,
			&teacher.Email,
			&teacher.Class,
			&teacher.Subject,
		)
	if err == sql.ErrNoRows {
		return models.Teacher{}, fmt.Errorf("teacher not found %v", err)
	} else if err != nil {
		return models.Teacher{}, fmt.Errorf("error quering the database %v", err)
	}
	return teacher, nil
}

func (t *Teachers) GetAllTeachers(firstName string, lastName string) (*sql.Rows, error) {
	query := `SELECT id, fisrt_name,last_name,email,class,subject FROM teachers WHERE 1=1`
	var args []any

	if firstName != "" {
		query += " AND first_name = ?"
		args = append(args, firstName)
	}

	if lastName != "" {
		query += " AND last_name = ?"
		args = append(args, lastName)
	}

	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query error %v", err)
	}
	defer rows.Close()
	return rows, nil
}
