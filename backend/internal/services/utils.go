package services

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Helper function to parse UUID
func parseUUID(s string) (pgtype.UUID, error) {
	if s == "" {
		return pgtype.UUID{}, nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: id, Valid: true}, nil
}
