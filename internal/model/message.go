package model

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type Message struct {
	ID              uuid.UUID      `json:"id,omitempty"`
	CreatedAt       time.Time      `json:"created_at,omitempty"`
	URL             string         `json:"url,omitempty"`
	RegexpPattern   sql.NullString `json:"regexp_pattern,omitempty"`
	RegexpMatch     sql.NullBool   `json:"regexp_match,omitempty"`
	StatusCode      sql.NullInt32  `json:"status_code,omitempty"`
	TimeToFirstByte *time.Duration `json:"time_to_first_byte,omitempty"`
}
