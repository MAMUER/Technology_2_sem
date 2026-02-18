package service

import "fmt"

type UserRecord struct {
	ID    int64
	Email string
	Role  string
	Hash  string
}

var ErrNotFound = fmt.Errorf("not found")

type UserRepo interface {
	ByEmail(email string) (UserRecord, error)
}
