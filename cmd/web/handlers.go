package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/labstack/echo/v4"
	"github.com/vulkan0n/superbchat/internal/models"
	"github.com/vulkan0n/superbchat/internal/validator"
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

func (app *application) getSettings(c echo.Context) error {
	return c.String(http.StatusOK, "Default")
}

func (app *application) postSettings(c echo.Context) error {
	return c.String(http.StatusOK, "Default")
}

func (app *application) postPay(c echo.Context) error {
	return c.String(http.StatusOK, "Default")
}

type userLoginForm struct {
	User                string `form:"user"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.User), "user", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusOK, "login.html", data)
		return
	}

	id, err := app.accounts.Authenticate(form.User, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Username or password incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusOK, "login.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authAccountId", id)

	http.Redirect(w, r, "/view", http.StatusSeeOther)
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
