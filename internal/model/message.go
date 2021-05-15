package model

import (
	"time"
)

type Message struct {
	CreatedAt       time.Time     `json:"created_at,omitempty"`
	URL             string        `json:"url,omitempty"`
	RegexpPattern   string        `json:"regexp_pattern,omitempty"`
	RegexpMatch     bool          `json:"regexp_match,omitempty"`
	StatusCode      int           `json:"status_code,omitempty"`
	TimeToFirstByte time.Duration `json:"time_to_first_byte,omitempty"`
}
