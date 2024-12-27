package types

import (
	"time"

	"github.com/google/uuid"
)

type ImageStatus int

const (
	ImageStatusFailed ImageStatus = iota
	ImageStatusPending
	ImageStatusCompleted
)

type Image struct {
	ID          int `bun:"id,pk,autoincrement"`
	UserID      uuid.UUID
	ImageStatus ImageStatus
	Prompt      string
	Deleted     bool      `bun:"default:'false'"`
	CreatedAt   time.Time `bun:"default:'now()'"`
	DeletedAt   time.Time
}
