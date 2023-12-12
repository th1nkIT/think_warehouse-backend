package utility

import "github.com/google/uuid"

func GenerateGoogleUUID() string {
	return uuid.New().String()
}
