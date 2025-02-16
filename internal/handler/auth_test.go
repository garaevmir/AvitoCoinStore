package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/tests/mocks"
)

func TestAuthHandler_Login(t *testing.T) {
	e := echo.New()
	userRepo := new(mocks.UserRepositoryMock)
	authHandler := NewAuthHandler(userRepo, "test-secret-key")

	t.Run("Invalid request", func(t *testing.T) {
		reqBody := map[string]int{
			"Username": 1,
			"Password": 1,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		err := authHandler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("Successful registration of a new user", func(t *testing.T) {
		reqBody := model.AuthRequest{
			Username: "testuser",
			Password: "testpass",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		userRepo.On("GetUserByUsername", mock.Anything, "testuser").
			Return((*model.User)(nil), nil).Once()

		userRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Username == "testuser" &&
				bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("testpass")) == nil
		})).Return(nil).Once()

		err := authHandler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response model.AuthResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NotEmpty(t, response.Token)

		userRepo.AssertExpectations(t)
	})

	t.Run("Invalid registration", func(t *testing.T) {
		reqBody := model.AuthRequest{
			Username: "",
			Password: "",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		err := authHandler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response model.AuthResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Empty(t, response.Token)

		userRepo.AssertExpectations(t)
	})

	t.Run("Getting by username error", func(t *testing.T) {
		reqBody := model.AuthRequest{
			Username: "testuser",
			Password: "testpass",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		userRepo.On("GetUserByUsername", mock.Anything, "testuser").
			Return((*model.User)(nil), model.ErrInternalError).Once()

		err := authHandler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response model.AuthResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Empty(t, response.Token)

		userRepo.AssertExpectations(t)
	})

	t.Run("Creating user error", func(t *testing.T) {
		reqBody := model.AuthRequest{
			Username: "testuser",
			Password: "testpass",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		userRepo.On("GetUserByUsername", mock.Anything, "testuser").
			Return((*model.User)(nil), nil).Once()

		userRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Username == "testuser" &&
				bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("testpass")) == nil
		})).Return(model.ErrCreateUser).Once()

		err := authHandler.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response model.AuthResponse
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Empty(t, response.Token)

		userRepo.AssertExpectations(t)
	})

	t.Run("Incorrect password", func(t *testing.T) {
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

		userRepo.AssertExpectations(t)
	})
}
