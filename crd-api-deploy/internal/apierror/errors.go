// Package apierror defines custom error types for the application.
package apierror

import "errors"

var (
	// ErrNotFound indicates that a requested resource was not found.
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput indicates that the user-provided input was invalid.
	ErrInvalidInput = errors.New("invalid input")

	// ErrK8sConflict indicates a conflict in the Kubernetes cluster, like a resource already existing.
	ErrK8sConflict = errors.New("kubernetes resource conflict")

	// ErrK8sAPIFailure indicates an unexpected error from the Kubernetes API.
	ErrK8sAPIFailure = errors.New("kubernetes api failure")
)
