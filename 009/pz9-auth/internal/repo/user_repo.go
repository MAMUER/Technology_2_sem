package repo

import (
	"context"
	"errors"
	"log"

	"example.com/pz9-auth/internal/core"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailTaken = errors.New("email already in use")

type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&core.User{})
}

func (r *UserRepo) Create(ctx context.Context, u *core.User) error {
	log.Printf("=== DEBUG: Starting Create for email: %s", u.Email)
	
	// Сначала проверяем существует ли пользователь с таким email
	var existing core.User
	err := r.db.WithContext(ctx).Where("email = ?", u.Email).First(&existing).Error
	
	log.Printf("=== DEBUG: Check existing user error: %v", err)
	log.Printf("=== DEBUG: Is ErrRecordNotFound? %v", errors.Is(err, gorm.ErrRecordNotFound))
	
	if err == nil {
		log.Printf("=== DEBUG: Email already taken")
		return ErrEmailTaken
	}
	
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("=== DEBUG: Unexpected DB error in check: %v, type: %T", err, err)
		return err
	}

	log.Printf("=== DEBUG: No existing user, creating new...")
	
	// Создаем нового пользователя
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		log.Printf("=== DEBUG: DB error in Create: %v", err)
		return err
	}
	
	log.Printf("=== DEBUG: User created successfully with ID: %d", u.ID)
	return nil
}

func (r *UserRepo) ByEmail(ctx context.Context, email string) (core.User, error) {
	var u core.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return core.User{}, ErrUserNotFound
	}
	return u, err
}