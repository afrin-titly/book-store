package infrastucture_test

import (
	"book-apis/domain"
	"book-apis/infrastucture"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepositoryDB_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %v", err)
	}
	defer db.Close()

	type testCase struct {
		name      string
		expected  *domain.User
		input     *domain.User
		mockSetup func()
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	tests := []testCase{
		{
			name: "Create user - success",
			expected: &domain.User{
				Name:    "Test User 1",
				Email:   "test@mail.com",
				Phone:   "07023456710",
				Address: "Setagaya, Tokyo",
			},
			input: &domain.User{
				Name:     "Test User 1",
				Email:    "test@mail.com",
				Password: string(hashedPassword),
				Phone:    "07023456710",
				Address:  "Setagaya, Tokyo",
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO users \\(name, email, password, phone, address\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
					WithArgs("Test User 1", "test@mail.com", string(hashedPassword), "07023456710", "Setagaya, Tokyo").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery("SELECT name, email, phone, address FROM users WHERE id = \\?").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"name", "email", "phone", "address"}).
						AddRow("Test User 1", "test@mail.com", "07023456710", "Setagaya, Tokyo"))
			},
		},
	}
	repo := infrastucture.NewUserRepositoryDB(db)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, err := repo.CreateUser(tc.input)
			if err != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
