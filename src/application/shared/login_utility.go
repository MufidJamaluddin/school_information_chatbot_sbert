package shared

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateClaims(userData *UserData, expirationTime time.Time) jwt.MapClaims {
	return jwt.MapClaims{
		"userName":    userData.UserName,
		"fullName":    userData.FullName,
		"roleGroupId": userData.RoleGroupId,
		"position":    userData.Position,
		"exp":         expirationTime.Unix(),
	}
}

func GetUserData(fiberCtx *fiber.Ctx) *UserData {
	user, ok := fiberCtx.Locals("user").(*jwt.Token)
	if !ok {
		return nil
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	userName, _ := GetStringFromMap(claims, "userName")
	roleGroupId, _ := GetUintFromMap(claims, "roleGroupId")
	fullName, _ := GetStringFromMap(claims, "fullName")
	position, _ := GetStringFromMap(claims, "position")

	return &UserData{
		UserName:    userName,
		RoleGroupId: roleGroupId,
		FullName:    fullName,
		Position:    position,
	}
}
