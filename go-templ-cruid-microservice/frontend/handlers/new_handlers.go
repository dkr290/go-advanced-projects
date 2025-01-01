package handlers

type Handlers struct {
	BackendService string
}

func NewHandlers(backendService string) *Handlers {
	return &Handlers{
		BackendService: backendService,
	}
}
