package config

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type AppConfig struct {
	Http                 HttpConf
	SqlDb                SqlDbConf
	WhatsAppConf         WhatsAppConf
	SecretKey            []byte
	PassExpirationInHour uint
	AllowOrigins         string
	AllowHeaders         string
	ExposeHeaders        string
	CookieKeyName        string
	CookieEncryptKey     string
	HeaderTokenKeyName   string
	IsLogged             bool
	IsFineTuneModel      bool
}

type HttpConf struct {
	Host           string
	Port           string
	Limiter        limiter.Config
	S2STimeoutSecs time.Duration
}

type WhatsAppConf struct {
	MessageAPIPattern string
	AccessToken       string
	VerifyToken       string
}

type SqlDbConf struct {
	Host                   string
	Username               string
	Password               string
	Name                   string
	Port                   string
	SSLMode                string
	IsLogged               bool
	MaxOpenConn            int
	MaxIdleConn            int
	MaxIdleTimeConnSeconds int64
	MaxLifeTimeConnSeconds int64
}

func New() *AppConfig {
	timeout, _ := strconv.ParseInt(os.Getenv("S2STimeoutSecs"), 10, 64)
	if timeout == 0 {
		timeout = 3000
	}

	httpConf := HttpConf{
		Host: os.Getenv("APP_HOST"),
		Port: os.Getenv("APP_PORT"),
		Limiter: limiter.Config{
			Next: func(c *fiber.Ctx) bool {
				return c.IP() == "127.0.0.1"
			},
			Max:        5,
			Expiration: 1 * time.Second,
		},
		S2STimeoutSecs: time.Duration(timeout) * time.Second,
	}

	sqlDbConf := SqlDbConf{
		IsLogged:               os.Getenv("IS_LOGGED") == "true",
		Host:                   os.Getenv("DB_HOST"),
		Username:               os.Getenv("DB_USERNAME"),
		Password:               os.Getenv("DB_PASSWORD"),
		Name:                   os.Getenv("DB_NAME"),
		Port:                   os.Getenv("DB_PORT"),
		SSLMode:                os.Getenv("DB_SSL_MODE"),
		MaxOpenConn:            10,
		MaxIdleConn:            1,
		MaxIdleTimeConnSeconds: 1 * 60 * 60,
		MaxLifeTimeConnSeconds: 5 * 60 * 60,
	}

	whatsAppConf := WhatsAppConf{
		MessageAPIPattern: os.Getenv("WHATSAPP_MSG_API_PATTERN"),
		AccessToken:       os.Getenv("WHATSAPP_TOKEN"),
		VerifyToken:       os.Getenv("VERIFY_TOKEN"),
	}

	passExpirationInHour, _ := strconv.ParseUint(os.Getenv("PASS_EXPIRATION_HOUR"), 10, 32)

	cookieKeyName := os.Getenv("COOKIE_KEY_NAME")
	headerTokenKeyName := os.Getenv("HEADER_TOKEN_KEY_NAME")
	cookieEncryptKey := os.Getenv("COOKIE_ENCRYPT_KEY")

	if cookieKeyName == "" {
		cookieKeyName = "TOKEN"
	}

	if headerTokenKeyName == "" {
		headerTokenKeyName = "TOKEN"
	}

	if cookieEncryptKey == "" {
		cookieEncryptKey = "286"
	}

	return &AppConfig{
		Http:                 httpConf,
		SqlDb:                sqlDbConf,
		WhatsAppConf:         whatsAppConf,
		SecretKey:            []byte(os.Getenv("APP_SECRET_KEY")),
		PassExpirationInHour: uint(passExpirationInHour),
		AllowOrigins:         os.Getenv("ALLOW_ORIGINS"),
		AllowHeaders:         os.Getenv("ALLOW_HEADERS"),
		ExposeHeaders:        os.Getenv("EXPOSE_HEADERS"),
		CookieKeyName:        cookieKeyName,
		CookieEncryptKey:     cookieEncryptKey,
		HeaderTokenKeyName:   headerTokenKeyName,
		IsLogged:             os.Getenv("IS_LOGGED") == "true",
		IsFineTuneModel:      os.Getenv("IS_FINE_TUNE_MODEL") == "true",
	}
}
