package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hathbanger/goDash/models"
	"github.com/labstack/echo"
)

func BulkUserController(c echo.Context) error {

	var u models.BulkUserCreate
	if err := c.Bind(&u); err != nil {
		fmt.Println("u", u.PhoneNumbers)
	}

	fmt.Println("u", u.PhoneNumbers)

	for _, v := range u.PhoneNumbers {
		user := models.NewUserModel(v, "", v, u.Organization)
		_ = user.Save()
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": "woooooo",
	})
}

func CreateUserController(c echo.Context) error {
	var u models.UserCreate
	if err := c.Bind(&u); err != nil {
		fmt.Println("u", u.Username)
	}

	var user *models.User
	var err error

	if len(u.Organization) == 0 {
		user = models.NewAdminUserModel(u.Username, u.Password, u.PhoneNumber, "")
		err = user.Save()
	} else {
		user = models.NewUserModel("", "", u.PhoneNumber, u.Organization)
		err = user.Save()
	}

	if err != nil {
		fmt.Println("SHIT!")
	}

	// TODO: need to make an admin user creation controller

	if err != nil {

		return c.JSON(
			http.StatusForbidden,
			"We're sorry! There's already a user with that username..")
	}

	userIdStr := user.Id.Hex()

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userIdStr
	claims["username"] = u.Username
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func LoginUserController(c echo.Context) error {
	var u models.UserLogin
	if err := c.Bind(&u); err != nil {
		fmt.Println("u", u.Username)
	}
	user, err := models.FindByUsernameModel(u.Username)
	if err != nil {
		return echo.ErrUnauthorized
	}

	if user.Password == u.Password {
		userIdStr := user.Id.Hex()
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = userIdStr
		claims["username"] = u.Username
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

func GetUserController(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)

	user, err := models.FindUserModel(userId)
	if err != nil {
		return err
	}

	fmt.Println(err)

	return c.JSON(http.StatusOK, user)
}

func UpdateUserController(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["id"].(string)
	username := claims["username"].(string)

	password := c.FormValue("password")
	models.UpdateUserModel(userId, username, password)
	user, err := models.FindUserModel(userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func RemoveUserController(c echo.Context) error {
	username := c.FormValue("username")
	user, err := models.FindByUsernameModel(username)
	err = models.DeleteUserModel(user.Id.Hex())
	if err != nil {
		return c.JSON(http.StatusNotFound, "not able to remove the account..")
	}

	return c.JSON(http.StatusOK, "User deleted!")
}
