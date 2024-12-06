package helpers

import "github.com/google/uuid"

func IsValidUUID(UUID string) bool {
	_, err := uuid.Parse(UUID)
	return err == nil
}
