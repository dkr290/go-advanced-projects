package handlers

import (
	"context"
	"sync"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

type TeacherHandlers struct {
	TeachersMap map[int]models.Teacher
	mutex       sync.Mutex
}

func NewTeachersHandler(teachers map[int]models.Teacher) *TeacherHandlers {
	return &TeacherHandlers{
		TeachersMap: teachers,
	}
}

func (h *TeacherHandlers) RootHandler(ctx context.Context, _ *struct{}) (*GreetingOutput, error) {
	resp := &GreetingOutput{}
	resp.Body.Message = "Hello from root Handler"
	return resp, nil
}

func (h *TeacherHandlers) TeachersGet(
	ctx context.Context,
	input *TeachersQueryInput,
) (*TeachersOutput, error) {
	response := TeachersOutput{}
	teacherList := make([]models.Teacher, 0, len(h.TeachersMap))
	for _, teacher := range h.TeachersMap {
		if (input.FirstName == "" || teacher.FirstName == input.FirstName) &&
			(input.LastName == "" || teacher.LastName == input.LastName) {
			teacherList = append(teacherList, teacher)
		}
	}
	response.Body.Status = "Sucess"
	response.Body.Count = len(teacherList)
	response.Body.Data = teacherList
	return &response, nil
}

func (h *TeacherHandlers) TeacherGet(ctx context.Context, input *struct {
	ID int `path:"id"`
},
) (*TeacherIDResponse, error) {
	resp := TeacherIDResponse{}

	teacher, exists := h.TeachersMap[input.ID]
	if !exists {
		return nil, huma.Error404NotFound("Teacher not found", nil)
	}
	resp.Body.Data = teacher
	return &resp, nil
}

func (h *TeacherHandlers) TeachersAdd(
	ctx context.Context,
	input *TeachersInput,
) (*TeachersOutput, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	addedTeachers := make([]models.Teacher, 0, len(input.Body.Teachers))

	maxID := 0
	for id := range h.TeachersMap {
		if id > maxID {
			maxID = id
		}
	}

	for _, newTeacher := range input.Body.Teachers {
		maxID++
		teacher := models.Teacher{
			ID:        maxID,
			FirstName: newTeacher.FirstName,
			LastName:  newTeacher.LastName,
			Class:     newTeacher.Class,
			Subject:   newTeacher.Subject,
		}
		h.TeachersMap[teacher.ID] = teacher
		addedTeachers = append(addedTeachers, teacher)
	}
	resp := &TeachersOutput{}
	resp.Body.Status = "Success"
	resp.Body.Count = len(addedTeachers)
	resp.Body.Data = addedTeachers
	return resp, nil
}
