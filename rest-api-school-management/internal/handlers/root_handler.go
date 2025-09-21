package handlers

import "context"

func (h *TeacherHandlers) RootHandler(ctx context.Context, _ *struct{}) (*GreetingOutput, error) {
	resp := &GreetingOutput{}
	resp.Body.Message = "Hello from root Handler"
	return resp, nil
}
