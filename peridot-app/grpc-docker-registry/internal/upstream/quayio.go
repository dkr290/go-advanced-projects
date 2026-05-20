package upstream

import "net/http"

type QuayRegistry struct {
	GenericRegistry
}

func NewQuayRegistry(client *http.Client) *GenericRegistry {
	return NewGenericRegistry("quay.io", client)
}

// NormalizeRepo — no prefix needed for quay.io
func (r *QuayRegistry) NormalizeRepo(repo string) string {
	return repo
}

