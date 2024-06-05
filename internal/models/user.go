package models

import "time"

type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	EncryptedPassword []byte    `json:"-"`
	BirthDate         time.Time `json:"birthdate"`
}
