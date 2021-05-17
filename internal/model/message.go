package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	ID              *uuid.UUID     `json:"id,omitempty"`
	CreatedAt       *time.Time     `json:"created_at,omitempty"`
	URL             *string        `json:"url,omitempty"`
	Error           *string        `json:"error,omitempty"`
	RegexpPattern   *string        `json:"regexp_pattern,omitempty"`
	RegexpMatch     *bool          `json:"regexp_match,omitempty"`
	StatusCode      *int32         `json:"status_code,omitempty"`
	TimeToFirstByte *time.Duration `json:"time_to_first_byte,omitempty"`
}
