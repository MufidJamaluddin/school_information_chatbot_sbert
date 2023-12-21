package login

import (
	"chatbot_be_go/src/application/login/dto"
	"chatbot_be_go/src/application/shared"
	appConf "chatbot_be_go/src/persistence/config"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
)

type ILoginHandler interface {
	LoginAction(c *fiber.Ctx) error
	LogoutAction(c *fiber.Ctx) error
}

type loginHandler struct {
	logger          *logrus.Logger
	appConfig       *appConf.AppConfig
	loginRepository ILoginRepository
}

var _ ILoginHandler = &loginHandler{}

func NewLoginHandler(
	logger *logrus.Logger,
	appConfig *appConf.AppConfig,
	loginRepository ILoginRepository,
) ILoginHandler {
	return &loginHandler{
		logger:          logger,
		appConfig:       appConfig,
		loginRepository: loginRepository,
	}
}

// LoginAction godoc
// @Tags "UC02 Admin Login"
// @Summary Login Admin & Get JWT Key
// @Description Get JWT Key for Authorization
// @Accept json
// @Produce json
// @Param q body dto.LoginDTO true "Login Data"
// @Success 202 {object} dto.AuthResponse
// @Failure 400 {object} string
// @Router /api/login [post]
func (l *loginHandler) LoginAction(c *fiber.Ctx) (err error) {
	var loginDto dto.LoginDTO
	var userData shared.UserData

	if err = json.Unmarshal(c.Body(), &loginDto); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if err = validator.Validate(&loginDto); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(
			fmt.Sprintf("Invalid Format: %v", err),
		)
	}

	if userData, err = l.loginRepository.GetUserData(c.Context(), loginDto.Username); err != nil {
		l.logger.Errorf("Error when Login with Username %s Pass %s: %v", loginDto.Username, loginDto.Password, err)
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Username or Password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userData.PasswordHash), []byte(loginDto.Password)); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Username or Password")
	}

	expiredAt := time.Now().Add(time.Hour * time.Duration(l.appConfig.PassExpirationInHour))
	claims := shared.CreateClaims(&userData, expiredAt)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(l.appConfig.SecretKey)
	if err != nil {
		l.logger.WithContext(c.Context()).WithError(err).Error("Error In Signed JWT Key")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     l.appConfig.CookieKeyName,
		Value:    signedToken,
		Expires:  expiredAt,
		HTTPOnly: true,
		SameSite: "strict",
	})

	return c.JSON(&dto.AuthResponse{
		Token:    signedToken,
		UserData: &userData,
	})
}

// LogoutAction godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC02 Admin Logout"
// @Summary Logout
// @Description Logout
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /api/login [delete]
func (l *loginHandler) LogoutAction(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     l.appConfig.CookieKeyName,
		Value:    "",
		HTTPOnly: true,
		SameSite: "strict",
	})

	c.ClearCookie(l.appConfig.CookieKeyName)

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Logout success!",
	})
}
