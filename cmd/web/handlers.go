package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/vulkan0n/superbchat/internal/models"
)

type PostCredentialsBody struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

func (app *application) postTest(c echo.Context) error {
	// Get the request
	r := c.Request()
	// Read the body
	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Body Request"})
	}
	n := PostCredentialsBody{
		Username: "default",
		Password: "default",
	}

	// equivalent of JSON.parse() in GO
	// By default Go passes arguments by value, meaning it creates a copy of the value, and a new pointer is created.
	// json.Unmarshall requires a reference (a pointer) to PostPersonBody and will update it internally.
	err = json.Unmarshal(b, &n)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}
	// Debug purpose
	app.infoLog.Println(n)
	// Update local instance (db...)

	return c.JSON(http.StatusOK, n)
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

func (app *application) getSettings(c echo.Context) error {
	return c.String(http.StatusOK, "Default")
}

func (app *application) postSettings(c echo.Context) error {
	return c.String(http.StatusOK, "Default")
}

func (app *application) postPay(c echo.Context) error {
	return c.String(http.StatusOK, "Default")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authAccountId")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type superbchatForm struct {
	Username     string
	AccountId    int
	NameMaxChars int
	MsgMaxChars  int
	MinAmnt      float64
	Checked      bool
}

func (app *application) superbchat(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "user")
	account, err := app.accounts.GetByUsername(username)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w, r)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Form = superbchatForm{
		Username:     account.Username,
		AccountId:    account.Id,
		NameMaxChars: account.NameMaxChars,
		MsgMaxChars:  account.MessageMaxChars,
		MinAmnt:      account.MinDonation,
		Checked:      account.IsDefaultShowAmount,
	}
	app.render(w, http.StatusOK, "superbchat.html", data)
}
