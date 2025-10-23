package repo

import (
	"errors"
	"golang.org/x/crypto/bcrypt"

	"example.com/pz10-auth/internal/core"
)

type UserRecord struct {
	ID    int64
	Email string
	Role  string
	Hash  []byte
}

type UserMem struct {
	users map[string]UserRecord
	usersByID map[int64]UserRecord
}

func NewUserMem() *UserMem {
	hash := func(s string) []byte {
		h, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
		return h
	}
	
	users := map[string]UserRecord{
		"admin@example.com": {ID: 1, Email: "admin@example.com", Role: "admin", Hash: hash("secret123")},
		"user@example.com":  {ID: 2, Email: "user@example.com", Role: "user", Hash: hash("secret123")},
		"user2@example.com": {ID: 3, Email: "user2@example.com", Role: "user", Hash: hash("secret123")},
	}
	
	usersByID := make(map[int64]UserRecord)
	for _, user := range users {
		usersByID[user.ID] = user
	}
	
	return &UserMem{
		users: users,
		usersByID: usersByID,
	}
}

var (
	ErrNotFound  = errors.New("user not found")
	ErrBadCreds = errors.New("bad credentials")
)

func (r *UserMem) CheckPassword(email, pass string) (*core.User, error) {
	userRecord, ok := r.users[email]
	if !ok {
		return nil, ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword(userRecord.Hash, []byte(pass)); err != nil {
		return nil, ErrBadCreds
	}

	user := &core.User{
		ID:    userRecord.ID,
		Email: userRecord.Email,
		Role:  userRecord.Role,
	}

	return user, nil
}

func (r *UserMem) GetUserByID(id int64) (*core.User, error) {
	userRecord, ok := r.usersByID[id]
	if !ok {
		return nil, ErrNotFound
	}

	user := &core.User{
		ID:    userRecord.ID,
		Email: userRecord.Email,
		Role:  userRecord.Role,
	}

	return user, nil
}