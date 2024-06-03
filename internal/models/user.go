package models

import "time"

type User struct {
	ID                int
	Email             string
	EncryptedPassword []byte
	BirthDate         time.Time
}
