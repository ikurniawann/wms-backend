// handlers/common.go
// Common handler utilities

package handlers

import (
	"time"

	"github.com/google/uuid"
)

// parseUUID parses string to uuid.UUID
func parseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

// parseDate parses date string
func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// parseDateTime parses datetime string
func parseDateTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}
