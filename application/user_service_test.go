package application_test

import (
	"book-apis/application"
	"book-apis/domain"
	"book-apis/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_CreateUser(t *testing.T) {
	repo := new(mocks.MockUserRepository)
	service := application.NewUserService(repo)
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	type testCase struct {
		name      string
		expected  *domain.User
		input     *domain.User
		mockSetup func()
	}

	tests := []testCase{
		{
			name: "Create user - valid input",
			expected: &domain.User{
				Name:    "Test User",
				Email:   "test@mail.com",
				Phone:   "07023456710",
				Address: "Setagaya, Tokyo",
			},
			input: &domain.User{
				Name:     "Test User",
				Email:    "test@mail.com",
				Password: "secret",
				Phone:    "07023456710",
				Address:  "Setagaya, Tokyo",
			},
			mockSetup: func() {
				repo.On("ExistsByEmail", "test@mail.com").Return(false, nil)
				repo.On("CreateUser", mock.MatchedBy(func(user *domain.User) bool {
					err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("secret"))
					return err == nil &&
						user.Name == "Test User" &&
						user.Email == "test@mail.com" &&
						user.Phone == "07023456710" &&
						user.Address == "Setagaya, Tokyo"
				})).Return(&domain.User{
					Name:    "Test User",
					Email:   "test@mail.com",
					Phone:   "07023456710",
					Address: "Setagaya, Tokyo",
				}, nil)
			},
		},

		{
			name:     "Create user - unique email failed",
			expected: nil,
			input: &domain.User{
				Name:     "Test User",
				Email:    "test@mail.com",
				Password: "secret",
				Phone:    "07023456710",
				Address:  "Setagaya, Tokyo",
			},
			mockSetup: func() {
				repo.On("ExistsByEmail", "test@mail.com").Return(true, nil)
				repo.On("CreateUser", mock.MatchedBy(func(user *domain.User) bool {
					err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("secret"))
					return err == nil &&
						user.Name == "Test User" &&
						user.Email == "test@mail.com" &&
						user.Phone == "07023456710" &&
						user.Address == "Setagaya, Tokyo"
				})).Return(nil, errors.New("Email already exixts"))
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo.Mock = mock.Mock{}
			tc.mockSetup()
			result, err := service.CreateUser(tc.input)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.expected, result)
				assert.NoError(t, err)
			}
		})
	}
}
