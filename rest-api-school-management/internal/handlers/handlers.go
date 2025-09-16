package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/internal/models"
)

type TeacherHandlers struct {
	TeachersMap map[int]models.Teacher
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
