package handlers

import (
	"context"
	"sync"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/dataops"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

type TeacherHandlers struct {
	mutex      sync.Mutex
	teachersDB dataops.DatabaseInf
}

func NewTeachersHandler(tdb dataops.DatabaseInf) *TeacherHandlers {
	return &TeacherHandlers{
		teachersDB: tdb,
	}
}

func (h *TeacherHandlers) TeacherGet(ctx context.Context, input *struct {
	ID int `path:"id"`
},
) (*TeacherIDResponse, error) {
	resp := TeacherIDResponse{}

	teacher, err := h.teachersDB.GetTeacherByID(input.ID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Error querying database:", err)
	}

	resp.Body.Data = teacher
	return &resp, nil
}

func (h *TeacherHandlers) TeachersGet(
	ctx context.Context,
	input *TeachersQueryInput,
) (*TeachersOutput, error) {
	response := TeachersOutput{}

	rows, err := h.teachersDB.GetAllTeachers(input.FirstName, input.LastName)
	if err != nil {
		return nil, huma.Error500InternalServerError("Error quering database", err)
	}

	teachersList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher
		rows.Scan(
			&teacher.ID,
			&teacher.FirstName,
			&teacher.LastName,
			&teacher.Email,
			&teacher.Class,
			&teacher.Subject,
		)
		if err != nil {
			return nil, huma.Error500InternalServerError("Error scanning database results", err)
		}
		teachersList = append(teachersList, teacher)
	}
	defer rows.Close()

	response.Body.Status = "Sucess"
	response.Body.Count = len(teachersList)
	response.Body.Data = teachersList
	return &response, nil
}

func (h *TeacherHandlers) TeachersAdd(
	ctx context.Context,
	input *TeachersInput,
) (*TeachersOutput, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	addedTeachers := make([]models.Teacher, len(input.Body.Teachers))

	for i, newTeacher := range input.Body.Teachers {

		teacher := models.Teacher{
			FirstName: newTeacher.FirstName,
			LastName:  newTeacher.LastName,
			Email:     newTeacher.Email,
			Class:     newTeacher.Class,
			Subject:   newTeacher.Subject,
		}
		id, err := h.teachersDB.InsertTeachers(&teacher)
		if err != nil {
			return nil, huma.Error500InternalServerError(
				"Error inserting data to the database",
				err,
			)
		}
		teacher.ID = int(id)
		addedTeachers[i] = teacher
	}

	resp := &TeachersOutput{}
	resp.Body.Status = "Success"
	resp.Body.Count = len(addedTeachers)
	resp.Body.Data = addedTeachers
	return resp, nil
}
