package dataops

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

type DatabaseInf interface {
	InsertTeachers(*models.Teacher) (int64, error)
	GetTeacherByID(int) (models.Teacher, error)
	GetAllTeachers(map[string]string, []string) (*sql.Rows, error)
	UpdateTeacher(int, models.Teacher) error
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

func (t *Teachers) GetAllTeachers(params map[string]string, sortBy []string) (*sql.Rows, error) {
	query := "SELECT id, first_name,last_name,email,class,subject FROM teachers WHERE 1=1"
	var args []any
	var orderByParts []string

	// Define a whitelist of allowed sortable columns
	allowedColumns := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	// filtering by map of params
	for param, dbField := range params {
		if dbField != "" {
			query += " AND " + param + " = ?"
			args = append(args, dbField)
		}
	}

	for _, criteria := range sortBy {
		parts := strings.Split(criteria, ":")
		if len(parts) == 2 {
			sortColumn := parts[0]
			sortOrder := strings.ToUpper(parts[1])

			if allowedColumns[sortColumn] && (sortOrder == "ASC" || sortOrder == "DESC") {
				orderByParts = append(orderByParts, fmt.Sprintf("%s %s", sortColumn, sortOrder))
			}
		}
	}

	if len(orderByParts) > 0 {
		query += " ORDER BY " + strings.Join(orderByParts, ", ")
	}
	fmt.Println(query)

	rows, err := t.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("database query error %v", err)
	}
	return rows, nil
}

func (t *Teachers) UpdateTeacher(id int, updatedTeacher models.Teacher) error {
	var existingTeacher models.Teacher

	row := t.db.QueryRow(
		"SELECT id ,first_name,last_name,email,class,subject from teachers WHERE id = ?",
		id,
	)
	err := row.Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Class,
		&existingTeacher.Subject,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return huma.Error500InternalServerError("Teacher not found", err)
		} else {
			return huma.NewError(http.StatusNotFound, "unable to retreive data", err)
		}
	}

	updatedTeacher.ID = existingTeacher.ID
	switch {
	}

	_, err = t.db.Exec(
		"UPDATE teachers SET first_name = ?, last_name = ? ,email = ? , class = ?,subject = ? WHERE id = ?  ",
		&updatedTeacher.FirstName,
		&updatedTeacher.LastName,
		&updatedTeacher.Email,
		&updatedTeacher.Class,
		&updatedTeacher.Subject,
		&updatedTeacher.ID,
	)
	if err != nil {
		return huma.Error500InternalServerError("Error updating teacher", err)
	}
	return nil
}
