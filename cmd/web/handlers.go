package main

import (
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/skip2/go-qrcode"
	"github.com/vulkan0n/superbchat/internal/fullstack"
	"github.com/vulkan0n/superbchat/internal/models"
	"github.com/vulkan0n/superbchat/internal/validator"
)

var AlertWidgetRefreshInterval = "10"

type csvLog struct {
	ID            string
	Name          string
	Message       string
	Amount        string
	DisplayToggle string
	Refresh       string
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "index.html", app.newTemplateData(r))
}

func (app *application) view(w http.ResponseWriter, r *http.Request) {
	accountId := app.sessionManager.GetInt(r.Context(), "authAccountId")

	superchats, err := app.superchats.GetFromAccount(accountId)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var paidSuperchats []*models.Superchat
	for _, superchat := range superchats {
		if superchat.IsPaid {
			paidSuperchats = append(paidSuperchats, superchat)
		}
	}

	data := app.newTemplateData(r)
	data.Superchats = paidSuperchats
	app.render(w, http.StatusOK, "view.html", data)
}

type settingsForm struct {
	Address             string  `form:"address"`
	MinDonation         float64 `form:"minDonation"`
	NameMaxChars        int     `form:"nameMaxChars"`
	MessageMaxChars     int     `form:"msgMaxChars"`
	IsDefaultShowAmount bool    `form:"showAmount"`
	validator.Validator `form:"-"`
}

func (app *application) settings(w http.ResponseWriter, r *http.Request) {
	accountId := app.sessionManager.GetInt(r.Context(), "authAccountId")
	account, err := app.accounts.Get(accountId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = settingsForm{
		Address:             account.Address,
		MinDonation:         account.MinDonation,
		NameMaxChars:        account.NameMaxChars,
		MessageMaxChars:     account.MessageMaxChars,
		IsDefaultShowAmount: account.IsDefaultShowAmount,
	}
	app.render(w, http.StatusOK, "settings.html", data)
}

func (app *application) settingsPost(w http.ResponseWriter, r *http.Request) {
	accountId := app.sessionManager.GetInt(r.Context(), "authAccountId")

	var form settingsForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.MaxChars(form.NameMaxChars, 100), "nameMaxChars", "Must be less than 100 characters long")
	form.CheckField(validator.MaxChars(form.MessageMaxChars, 1000), "msgMaxChars", "Must be less than 1000 characters long")
	form.CheckField(validator.Matches(form.Address, validator.AddressRX), "address", "Invalid address format")

	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusOK, "settings.html", data)
		return
	}

	account, err := app.accounts.Get(accountId)
	if err != nil {
		app.serverError(w, err)
		return
	}
	account.Address = form.Address
	account.MinDonation = form.MinDonation
	account.NameMaxChars = form.NameMaxChars
	account.MessageMaxChars = form.MessageMaxChars
	account.IsDefaultShowAmount = form.IsDefaultShowAmount
	err = app.accounts.Update(account)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Settings updated successfully")
	http.Redirect(w, r, "/view", http.StatusSeeOther)
}

func remove(stringSlice []string, stringToRemove string) []string {
	for i, v := range stringSlice {
		if v == stringToRemove {
			return append(stringSlice[:i], stringSlice[i+1:]...)
		}
	}
	return stringSlice
}

func (app *application) alertHandler(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("user")
	var v csvLog
	v.Refresh = AlertWidgetRefreshInterval
	auth := r.URL.Query().Get("auth")
	if auth == password {

		csvFile, err := os.Open("./cmd/log/alertqueue.csv")
		if err != nil {
			fmt.Println(err)
		}

		csvLines, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			fmt.Println(err)
		}
		defer func(csvFile *os.File) {
			err := csvFile.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(csvFile)

		// Remove top line of CSV file after displaying it
		if csvLines != nil {
			popFile, _ := os.OpenFile("./cmd/log/alertqueue.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			popFirst := csvLines[1:]
			w := csv.NewWriter(popFile)
			err := w.WriteAll(popFirst)
			if err != nil {
				fmt.Println(err)
			}
			defer func(popFile *os.File) {
				err := popFile.Close()
				if err != nil {
					fmt.Println(err)
				}
			}(popFile)
			v.ID = csvLines[0][0]
			v.Name = csvLines[0][1]
			v.Message = csvLines[0][2]
			v.Amount = csvLines[0][3]
			v.DisplayToggle = ""
		} else {
			v.DisplayToggle = "display: none;"
		}
	} else {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return // return http 401 unauthorized error
	}
	//files := []string{
	//	"./ui/html/pages/alert.html",
	//}

	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.errorLog.Fatal(err.Error())
	//	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//}

	//err = ts.Execute(w, v)
	//if err != nil {
	//	app.errorLog.Fatal(err.Error())
	//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//	}
}

type payPostForm struct {
	AccountId  int     `form:"accountId"`
	Name       string  `form:"name"`
	Amount     float64 `form:"amount"`
	Message    string  `form:"message"`
	ShowAmount bool    `form:"showAmount"`
}

type payForm struct {
	Amount    float64
	Name      string
	Message   string
	Address   string
	AddressQR string
	CheckURL  string
}

func (app *application) payPost(w http.ResponseWriter, r *http.Request) {
	var postForm payPostForm
	err := app.decodePostForm(r, &postForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	account, err := app.accounts.Get(postForm.AccountId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if postForm.Amount == 0 {
		postForm.Amount = account.MinDonation
	}
	if postForm.Name == "" {
		postForm.Name = "Anonymous"
	}

	superchatId, err := app.superchats.Insert("", postForm.Name, postForm.Message, postForm.Amount,
		!postForm.ShowAmount, postForm.AccountId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	tmp, _ := qrcode.Encode(fmt.Sprintf("%s?amount=%f", account.Address, postForm.Amount), qrcode.Low, 320)

	form := payForm{
		Amount:    postForm.Amount,
		Name:      postForm.Name,
		Message:   postForm.Message,
		Address:   account.Address,
		AddressQR: base64.StdEncoding.EncodeToString(tmp),
		CheckURL:  fmt.Sprintf("%v/%v", superchatId, postForm.AccountId),
	}
	data := app.newTemplateData(r)
	data.Form = form

	app.render(w, http.StatusOK, "pay.html", data)
}

type checkForm struct {
	Receipt     string
	Autorefresh bool
}

func (app *application) check(w http.ResponseWriter, r *http.Request) {
	superchatId, _ := strconv.Atoi(chi.URLParam(r, "superchatId"))
	accountId, _ := strconv.Atoi(chi.URLParam(r, "accountId"))
	accountSuperchats, err := app.superchats.GetFromAccount(accountId)
	if err != nil {
		app.serverError(w, err)
		return
	}
	account, err := app.accounts.Get(accountId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form := checkForm{
		Receipt:     "Waiting for payment...",
		Autorefresh: true,
	}
	txsWallet, err := fullstack.GetTXs(account.Address)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var currentSuperchat models.Superchat
	for _, superchat := range accountSuperchats {
		if superchat.Id == superchatId {
			currentSuperchat = *superchat
		}
		if superchat.IsPaid {
			txsWallet = remove(txsWallet, superchat.TxId)
		}
	}

	txsBatchSize := 20

	for i := 0; i < len(txsWallet); i += txsBatchSize {
		j := i + txsBatchSize
		if j > len(txsWallet) {
			j = len(txsWallet)
		}
		txsBatch := txsWallet[i:j]
		txsDetailsResp, err := fullstack.GetTxsDetailsResponse(txsBatch)
		if err != nil {
			app.serverError(w, err)
			return
		}
		for _, tx := range txsDetailsResp.Transactions {
			for _, vout := range tx.Details.Vout {
				if vout.Value == currentSuperchat.Amount {
					app.superchats.SetAsPaid(tx.TxId, currentSuperchat.Id)
					form.Autorefresh = false
					form.Receipt = fmt.Sprintf("%f BCH Received! Superchat sent", currentSuperchat.Amount)
				}
			}
		}

	}
	data := app.newTemplateData(r)
	data.CustomStyle = "style-check.css"
	data.Autorefresh = form.Autorefresh
	data.Form = form
	app.render(w, http.StatusOK, "check.html", data)
}

type userSignupForm struct {
	User                string `form:"user"`
	Password            string `form:"password"`
	RepeatedPassword    string `form:"repeated-password"`
	Address             string `form:"address"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.User), "user", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Must be at least 8 characters long")
	form.CheckField(validator.EqualValue(form.Password, form.RepeatedPassword), "repeatedPassword", "Password doesn't match")
	form.CheckField(validator.Matches(form.Address, validator.AddressRX), "address", "Invalid address format")

	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusOK, "signup.html", data)
		return
	}
	err = app.accounts.Insert(form.User, form.Password, form.Address)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateUser) {
			form.AddFieldError("user", "Username already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusOK, "signup.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
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
