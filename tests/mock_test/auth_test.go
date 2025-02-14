package mock_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garaevmir/avitocoinstore/internal/handler"
	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/tests/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthHandler_Login(t *testing.T) {
	e := echo.New()
	userRepo := new(mocks.UserRepositoryMock)
	authHandler := handler.NewAuthHandler(userRepo, "test-secret-key")

	t.Run("Successful registration of a new user", func(t *testing.T) {
		reqBody := model.AuthRequest{
			Username: "new_user",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		userRepo.On("GetUserByUsername", mock.Anything, "new_user").
			Return((*model.User)(nil), nil).Once()

		userRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Username == "new_user" &&
				bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("password123")) == nil
		})).Return(nil).Once()

		err := authHandler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response model.AuthResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NotEmpty(t, response.Token)

		userRepo.AssertExpectations(t)
	})

	t.Run("Error: incorrect password", func(t *testing.T) {
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte("correct_pass"), bcrypt.DefaultCost)
		existingUser := &model.User{
			Username:     "existing_user",
			PasswordHash: string(hashedPass),
		}

		reqBody := model.AuthRequest{
			Username: "existing_user",
			Password: "wrong_pass",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		userRepo.On("GetUserByUsername", mock.Anything, "existing_user").
			Return(existingUser, nil).Once()

		err := authHandler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var errorResp model.ErrorResponse
		json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.Equal(t, "invalid credentials", errorResp.Errors)

		userRepo.AssertExpectations(t)
	})
}
