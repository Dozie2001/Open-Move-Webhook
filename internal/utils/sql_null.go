package utils

import (
	"database/sql"
	"time"
)

func SQLNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

func SQLNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: true}
}
