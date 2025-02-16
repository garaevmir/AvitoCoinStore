package handler

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/garaevmir/avitocoinstore/internal/repository"
)

// A structure for a authentication handler
type AuthHandler struct {
	userRepo repository.UserRepositoryInt
	secret   string
}

// Constructor for authentication handler
func NewAuthHandler(userRepo repository.UserRepositoryInt, secret string) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, secret: secret}
}

// Function for /api/auth request
func (h *AuthHandler) Login(c echo.Context) error {
	var req model.AuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrInvalidRequest)
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, model.ErrInvalidCredentials)
	}

	user, err := h.userRepo.GetUserByUsername(c.Request().Context(), req.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrInternalError)
	}

	if user == nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		newUser := &model.User{
			Username:     req.Username,
			PasswordHash: string(hashedPassword),
			Coins:        1000,
		}
		if err := h.userRepo.CreateUser(c.Request().Context(), newUser); err != nil {
			return c.JSON(http.StatusInternalServerError, model.ErrCreateUser)
		}
		user = newUser
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, model.ErrInvalidCredentials)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": user.ID})
	tokenString, _ := token.SignedString([]byte(h.secret))
	return c.JSON(http.StatusOK, model.AuthResponse{Token: tokenString})
}
