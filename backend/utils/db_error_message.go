package utils

import (
	"regexp"

	"github.com/jackc/pgx/v5/pgconn"
)

func ParseDBError(err error) (status int, payload interface{}) {
	if err == nil {
		return 200, nil
	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "23505": 
			field := extractFieldFromDetail(pgErr.Detail)
			if field == "" {
				field = pgErr.ConstraintName
			}
			msg := field + " already exists"
			return 400, map[string]string{field: msg}

		case "23503":
			field := extractFieldFromDetail(pgErr.Detail)
			if field == "" {
				field = pgErr.ConstraintName
			}
			msg := field + " is invalid or does not exist"
			return 400, map[string]string{field: msg}

		case "23502":
			field := pgErr.ColumnName
			if field == "" {
				field = extractFieldFromDetail(pgErr.Detail)
			}
			if field == "" {
				field = "field"
			}
			msg := field + " cannot be null"
			return 400, map[string]string{field: msg}

		case "23514": 
			return 400, map[string]string{"error": pgErr.Message}

		default:
			return 500, map[string]string{"error": pgErr.Message}
		}
	}

	return 500, map[string]string{"error": err.Error()}
}

func extractFieldFromDetail(detail string) string {
	re := regexp.MustCompile(`\(([^)]+)\)=`)
	m := re.FindStringSubmatch(detail)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}
