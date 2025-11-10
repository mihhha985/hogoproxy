package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"test/internal/responder"

	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	tokenAuth *jwtauth.JWTAuth
	user      *User
}

func NewAuthController(tokenAuth *jwtauth.JWTAuth, user *User) *AuthController {
	return &AuthController{
		tokenAuth: tokenAuth,
		user:      user,
	}
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя и возвращает JWT токен для аутентификации
// @Tags auth
// @Accept json
// @Produce json
// @Param request body User true "Данные пользователя для регистрации"
// @Success 200 {object} TokenResponse "JWT токен успешно создан"
// @Failure 400 {object} ErrorResponse "Некорректный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /register [post]
func (c *AuthController) Register() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var data User
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			responder.ErrorBadRequest(w, err)
			return
		}

		_, tokenString, err := c.tokenAuth.Encode(map[string]interface{}{"email": data.Username})
		if err != nil {
			responder.ErrorInternal(w, err)
			return
		}

		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			responder.ErrorInternal(w, err)
			return
		}

		c.user.Password = string(hashedBytes)
		c.user.Username = data.Username
		responder.OutputJSON(w, TokenResponse{Token: tokenString})
	}
}

// Login godoc
// @Summary Вход пользователя
// @Description Аутентифицирует пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body User true "Учетные данные пользователя"
// @Success 200 {object} TokenResponse "JWT токен успешно создан"
// @Failure 400 {object} ErrorResponse "Некорректный запрос"
// @Failure 401 {object} ErrorResponse "Неверное имя пользователя или пароль"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /login [post]
func (c *AuthController) Login() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var data User
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			responder.ErrorBadRequest(w, err)
			return
		}

		_, tokenString, err := c.tokenAuth.Encode(TokenClaims{"username": data.Username})
		if err != nil {
			responder.ErrorInternal(w, err)
			return
		}

		if data.Username != c.user.Username {
			responder.ErrorUnauthorized(w, errors.New("invalid username or password"))
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(c.user.Password), []byte(data.Password)); err != nil {
			responder.ErrorUnauthorized(w, errors.New("invalid username or password"))
			return
		}

		responder.OutputJSON(w, TokenResponse{Token: tokenString})
	}
}
