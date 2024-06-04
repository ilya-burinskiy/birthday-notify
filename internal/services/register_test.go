package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ilya-burinskiy/birthday-notify/internal/auth"
	"github.com/ilya-burinskiy/birthday-notify/internal/models"
	"github.com/ilya-burinskiy/birthday-notify/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type userCreator struct{ mock.Mock }

func (s *userCreator) CreateUser(
	ctx context.Context,
	login string,
	encryptedPassword []byte,
	birthdayDate time.Time,
) (models.User, error) {

	args := s.Called(ctx, login, encryptedPassword)
	return args.Get(0).(models.User), args.Error(1)
}

func TestRegister(t *testing.T) {
	type expected struct {
		jwtStr string
		err    error
	}
	type createUserResult struct {
		user models.User
		err  error
	}

	usrCreator := new(userCreator)
	registerSrv := services.NewRegisterService(usrCreator)
	testCases := []struct {
		name         string
		login        string
		password     string
		birthdayDate time.Time
		createRes    createUserResult
		expected     expected
	}{
		{
			name:         "returns JWT string",
			login:        "login",
			password:     "password",
			birthdayDate: time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
			createRes: createUserResult{
				user: models.User{ID: 1},
			},
			expected: expected{
				jwtStr: buildJWTString(t, 1),
			},
		},
		{
			name:         "returns error if failed to create user",
			login:        "login",
			password:     "password",
			birthdayDate: time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
			createRes: createUserResult{
				err: errors.New("error"),
			},
			expected: expected{
				err: fmt.Errorf("failed to register user: %w", errors.New("error")),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.TODO()
			createCall := usrCreator.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
				Return(tc.createRes.user, tc.createRes.err)
			defer createCall.Unset()

			jwtStr, err := registerSrv.Register(ctx, tc.login, tc.password, tc.birthdayDate)
			if err == nil {
				assert.Equal(
					t,
					userIDFromJWT(t, tc.expected.jwtStr),
					userIDFromJWT(t, jwtStr),
				)
			} else {
				assert.EqualError(t, err, "failed to register user: error")
			}
		})
	}
}

func buildJWTString(t *testing.T, userID int) string {
	jwtStr, error := auth.BuildJWTString(userID)
	require.NoError(t, error)

	return jwtStr
}
