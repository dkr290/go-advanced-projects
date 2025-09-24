package handlers

import (
	"context"
	"net/http"
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

	params := map[string]string{
		"first_name": input.FirstName,
		"last_name":  input.LastName,
		"email":      input.Email,
		"class":      input.Class,
		"subject":    input.Subject,
	}

	sortBy := input.SortBy
	// filtering by params basically with query parameters anf filtering
	rows, err := h.teachersDB.GetAllTeachers(params, sortBy)
	if err != nil {
		return nil, huma.Error500InternalServerError("Error quering database", err)
	}

	teachersList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(
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

func (h *TeacherHandlers) UpdateTeacherHandler(
	ctx context.Context,
	input *TeachersUpdateInput,
) (*struct{}, error) {
	id := input.Body.Teacher.ID
	if id <= 0 {
		return nil, huma.NewError(http.StatusBadRequest, "invalid teacher id", nil)
	}

	teacher := models.Teacher{
		ID:        input.Body.Teacher.ID,
		FirstName: input.Body.Teacher.FirstName,
		LastName:  input.Body.Teacher.LastName,
		Email:     input.Body.Teacher.Email,
		Class:     input.Body.Teacher.Class,
		Subject:   input.Body.Teacher.Subject,
	}

	err := h.teachersDB.UpdateTeacher(input.Body.Teacher.ID, teacher)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
