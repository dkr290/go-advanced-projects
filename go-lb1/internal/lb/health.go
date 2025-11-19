package lb

import (
	"sync"
	"time"
)

type Backend struct {
	Address         string
	Healthy         bool
	QuarantineUntil time.Time // Time when backend can be re-evaluated
	mutex           sync.RWMutex
}
