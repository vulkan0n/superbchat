package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/vulkan0n/superbchat/internal/models"
)

type PostCredentialsBody struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

func (app *application) postUserSignup(c echo.Context) error {
	r := c.Request()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Body Request"})
	}
	n := PostCredentialsBody{
		Username: "default",
		Password: "default",
	}
	err = json.Unmarshal(b, &n)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}
	err = app.accounts.Insert(n.Username, n.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUser) {
			app.errorLog.Println(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username already in use"})
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}
	return c.String(http.StatusOK, "User added")
}

type JwtClaims struct {
	UserId int `json:"userId"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("BCH_is_awsome") // Secret key to sign tokens

func (app *application) postUserLogin(c echo.Context) error {
	r := c.Request()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Body Request"})
	}
	credential := PostCredentialsBody{
		Username: "default",
		Password: "default",
	}
	err = json.Unmarshal(b, &credential)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	id, err := app.accounts.Authenticate(credential.Username, credential.Password)

	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {

			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid username or password"})
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}
	expirationTime := time.Now().Add(12 * time.Hour)
	claims := &JwtClaims{
		UserId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Token generation failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token":  tokenString,
		"userId": strconv.Itoa(id),
	})
}

type PostTokenBody struct {
	Token string `json:"token"`
}

func (app *application) postVerifyToken(c echo.Context) error {
	r := c.Request()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Body Request"})
	}
	tokenBody := PostTokenBody{
		Token: "default",
	}
	err = json.Unmarshal(b, &tokenBody)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	userId, err := validateToken(tokenBody.Token)

	if userId >= 0 {
		return c.JSON(http.StatusOK, map[string]int{"userId": userId})
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid token"})
	}
}

func validateToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return -1, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims.UserId, nil
	} else {
		return -1, nil
	}
}

type UserIdBody struct {
	UserId int `json:"userId"`
}

type UserInfoResponse struct {
	UserId     int    `json:"userId"`
	Username   string `json:"username"`
	Address    string `json:"address"`
	TknAddress string `json:"tknAddress"`
	ShowAmount bool   `json:"showAmount"`
}

func (app *application) getUserInfo(c echo.Context) error {
	user := c.Param("user")

	userInfo, err := app.accounts.GetByUsername(user)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	response := UserInfoResponse{
		UserId:     userInfo.Id,
		Username:   userInfo.Username,
		Address:    userInfo.Address,
		TknAddress: userInfo.TknAddress,
		ShowAmount: userInfo.IsDefaultShowAmount,
	}

	return c.JSON(http.StatusOK, response)
}

type postSuperbchatBody struct {
	Name      string  `json:"name"`
	Message   string  `json:"message"`
	Amount    float64 `json:"amount"`
	IsHidden  bool    `json:"isHidden"`
	Recipient int     `json:"recipient"`
	IsTkn     bool    `json:"isTkn"`
}

func (app *application) postSuperbchat(c echo.Context) error {
	r := c.Request()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Body Request"})
	}
	n := postSuperbchatBody{
		Name:      "default",
		Message:   "default",
		Amount:    0,
		IsHidden:  false,
		Recipient: 0,
		IsTkn:     false,
	}
	err = json.Unmarshal(b, &n)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}
	_, err = app.superchats.Insert("", n.Name, n.Message, n.Amount, "", 0, n.IsHidden, n.Recipient)

	if err != nil {
		app.infoLog.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.String(http.StatusOK, "Superbchat received")
}

func (app *application) getSettings(c echo.Context) error {
	r := c.Request()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Body Request"})
	}
	credential := UserIdBody{
		UserId: 0,
	}
	err = json.Unmarshal(b, &credential)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	return c.String(http.StatusOK, "Default")
}

func (app *application) postSettings(c echo.Context) error {
	return c.String(http.StatusOK, "Default")
}
