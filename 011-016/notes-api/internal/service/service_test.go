package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubRepo struct {
	users map[string]UserRecord
}

func (r stubRepo) ByEmail(email string) (UserRecord, error) {
	u, ok := r.users[email]
	if !ok {
		return UserRecord{}, ErrNotFound
	}
	return u, nil
}

// Вспомогательная функция для хеширования
func hash(password string) string {
	return "hashed_" + password
}

func TestService_FindIDByEmail(t *testing.T) {
	// Подготовка тестовых данных
	users := map[string]UserRecord{
		"admin@example.com": {ID: 1, Email: "admin@example.com", Role: "admin", Hash: hash("secret123")},
		"user@example.com":  {ID: 2, Email: "user@example.com", Role: "user", Hash: hash("secret123")},
		"user2@example.com": {ID: 3, Email: "user2@example.com", Role: "user", Hash: hash("secret123")},
	}

	repo := stubRepo{users: users}
	service := New(repo)

	t.Run("AdminFound", func(t *testing.T) {
		// «найден» - админ
		email := "admin@example.com"
		expectedID := int64(1)

		id, err := service.FindIDByEmail(email)

		require.NoError(t, err)
		assert.Equal(t, expectedID, id)
	})

	t.Run("UserFound", func(t *testing.T) {
		// «найден» - обычный пользователь
		email := "user@example.com"
		expectedID := int64(2)

		id, err := service.FindIDByEmail(email)

		require.NoError(t, err)
		assert.Equal(t, expectedID, id)
	})

	t.Run("User2Found", func(t *testing.T) {
		// «найден» - второй пользователь
		email := "user2@example.com"
		expectedID := int64(3)

		id, err := service.FindIDByEmail(email)

		require.NoError(t, err)
		assert.Equal(t, expectedID, id)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// «не найден»
		email := "nonexistent@example.com"

		id, err := service.FindIDByEmail(email)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Equal(t, int64(0), id)
	})

	t.Run("EmptyEmail", func(t *testing.T) {
		// пустой email
		email := ""

		id, err := service.FindIDByEmail(email)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Equal(t, int64(0), id)
	})
}

// Дополнительные тесты
func TestService_FindIDByEmail_TableDriven(t *testing.T) {
	users := map[string]UserRecord{
		"admin@example.com": {ID: 1, Email: "admin@example.com", Role: "admin", Hash: hash("secret123")},
		"user@example.com":  {ID: 2, Email: "user@example.com", Role: "user", Hash: hash("secret123")},
		"user2@example.com": {ID: 3, Email: "user2@example.com", Role: "user", Hash: hash("secret123")},
	}

	repo := stubRepo{users: users}
	service := New(repo)

	testCases := []struct {
		name        string
		email       string
		expectedID  int64
		expectError bool
	}{
		{"Found_Admin", "admin@example.com", 1, false},
		{"Found_User", "user@example.com", 2, false},
		{"Found_User2", "user2@example.com", 3, false},
		{"NotFound_Unknown", "unknown@example.com", 0, true},
		{"NotFound_Empty", "", 0, true},
		{"NotFound_DifferentDomain", "admin@gmail.com", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := service.FindIDByEmail(tc.email)

			if tc.expectError {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrNotFound)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.expectedID, id)
		})
	}
}
