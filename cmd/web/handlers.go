package main

import (
	"bufio"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/go-chi/chi"
	"github.com/skip2/go-qrcode"
	"github.com/vulkan0n/superbchat/internal/models"
	"github.com/vulkan0n/superbchat/internal/validator"
)

var AlertWidgetRefreshInterval = "10"

type checkPage struct {
	Addy     string
	PayID    string
	Received float64
	Meta     string
	Name     string
	Msg      string
	Receipt  string
}

type superChat struct {
	Name     string
	Message  string
	Amount   string
	Address  string
	QRB64    string
	PayID    string
	CheckURL string
}

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

func (app *application) checkHandler(w http.ResponseWriter, r *http.Request) {
	account, err := app.accounts.Get(app.sessionManager.GetInt(r.Context(), "authAccountId"))
	if err != nil {
		app.serverError(w, err)
	}

	var c checkPage
	c.Meta = `<meta http-equiv="Refresh" content="5">`
	c.Addy = account.Address
	c.Received, _ = strconv.ParseFloat(r.FormValue("amount"), 64)
	c.Name = truncateStrings(r.FormValue("name"), account.NameMaxChars)
	c.Msg = truncateStrings(r.FormValue("msg"), account.MessageMaxChars)
	c.Receipt = "Waiting for payment..."

	var txsWallet []string
	getTXs(account.Address, &txsWallet)
	var txsPaidLog []string
	getPaidLogTxs(&txsPaidLog)
	for _, txToRemove := range txsPaidLog {
		txsWallet = remove(txsWallet, txToRemove)
	}

	txsBatchSize := 20

	for i := 0; i < len(txsWallet); i += txsBatchSize {
		j := i + txsBatchSize
		if j > len(txsWallet) {
			j = len(txsWallet)
		}
		txsBatch := txsWallet[i:j]
		txsDetailsResp := &transactionsDetailsResponse{}
		getTxsDetailsResponse(txsDetailsResp, txsBatch)

		for _, tx := range txsDetailsResp.Transactions {
			for _, vout := range tx.Details.Vout {
				if vout.Value == c.Received {
					appendTxToLog(tx.TxId)

					c.Meta = ""
					setCheckReceipt(&c.Receipt, c.Received)

					if c.Msg == "" {
						c.Msg = "⠀"
					}
					c.PayID = tx.TxId
					if c.Received >= account.MinDonation {
						appendTxToCSVs(c.PayID, c.Name, c.Msg, c.Received, r.FormValue("show"))
					}
				}
			}
		}

	}
	files := []string{
		"./ui/html/pages/check.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", c)
	if err != nil {
		app.errorLog.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func condenseSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
func truncateStrings(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for !utf8.ValidString(s[:n]) {
		n--
	}
	return s[:n]
}
func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

func appendTxToCSVs(cPayID string, cName string, cMsg string, cReceived float64, show string) {
	f, err := os.OpenFile("./cmd/log/superchats.csv",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	csvAppend := fmt.Sprintf(`"%s","%s","%s","%s"`, cPayID, html.EscapeString(cName), html.EscapeString(cMsg), fmt.Sprint(cReceived))
	if show != "true" {
		csvAppend = fmt.Sprintf(`"%s","%s","%s","%s (hidden)"`, cPayID, html.EscapeString(cName), html.EscapeString(cMsg), fmt.Sprint(cReceived))
	}
	a, err := os.OpenFile("./cmd/log/alertqueue.csv",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer func(a *os.File) {
		err := a.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(a)
	fmt.Println(csvAppend)

	if _, err := f.WriteString(csvAppend + "\n"); err != nil {
		log.Println(err)
	}

	if show != "true" {
		csvAppend = fmt.Sprintf(`"%s","%s","%s","???"`, cPayID, html.EscapeString(cName), html.EscapeString(cMsg))
	}

	if _, err := a.WriteString(csvAppend + "\n"); err != nil {
		log.Println(err)
	}
}

func setCheckReceipt(receiptPtr *string, received float64) {
	*receiptPtr = "<b>" + fmt.Sprint(received) + " BCH Received! Superchat sent</b>"
}

func getPaidLogTxs(txsPaidLog *[]string) {
	file, err := os.Open("./cmd/log/paid.log")
	if err != nil {
		log.Fatalf("failed to open ")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		*txsPaidLog = append(*txsPaidLog, scanner.Text())
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func appendTxToLog(txId string) {
	f, err := os.OpenFile("./cmd/log/paid.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	if _, err := f.WriteString(txId + "\n"); err != nil {
		log.Println(err)
	}
}

func remove(stringSlice []string, stringToRemove string) []string {
	for i, v := range stringSlice {
		if v == stringToRemove {
			return append(stringSlice[:i], stringSlice[i+1:]...)
		}
	}
	return stringSlice
}

func (app *application) topwidgetHandler(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()
	if !ok {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if (u == "test") && (p == "tust") {
		csvFile, err := os.Open("./cmd/log/superchats.csv")
		if err != nil {
			fmt.Println(err)
		}
		defer func(csvFile *os.File) {
			err := csvFile.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(csvFile)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return // return http 401 unauthorized error
	}
	files := []string{
		"./ui/html/pages/top.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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
	files := []string{
		"./ui/html/pages/alert.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = ts.Execute(w, v)
	if err != nil {
		app.errorLog.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

type payPostForm struct {
	AccountId  int     `form:"accountId"`
	Name       string  `form:"name"`
	Amount     float64 `form:"amount"`
	Message    string  `form:"message"`
	ShowAmount bool    `form:"showAmount"`
}

func (app *application) payPost(w http.ResponseWriter, r *http.Request) {
	var form payPostForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	account, err := app.accounts.Get(form.AccountId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if form.Amount == 0 {
		form.Amount = account.MinDonation
	}
	if form.Name == "" {
		form.Name = "Anonymous"
	}

	err = app.superchats.Insert("", form.Name, form.Message, form.Amount,
		!form.ShowAmount, form.AccountId)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/"+account.Username, http.StatusSeeOther)
}

func (app *application) pay(w http.ResponseWriter, r *http.Request) {
	accountName := r.URL.Query().Get("user")
	account, err := app.accounts.GetByUsername(accountName)
	if err != nil {
		app.serverError(w, err)
		return
	}
	var s superChat
	s.Amount = html.EscapeString(r.FormValue("amount"))
	if r.FormValue("amount") == "" {
		s.Amount = fmt.Sprint(account.MinDonation)
	}
	if r.FormValue("name") == "" {
		s.Name = "Anonymous"
	} else {
		s.Name = html.EscapeString(truncateStrings(condenseSpaces(r.FormValue("name")), account.NameMaxChars))
	}
	s.Message = html.EscapeString(truncateStrings(condenseSpaces(r.FormValue("message")), account.MessageMaxChars))
	s.Address = account.Address

	params := url.Values{}
	params.Add("amount", s.Amount)
	params.Add("name", s.Name)
	params.Add("msg", r.FormValue("message"))
	params.Add("show", html.EscapeString(r.FormValue("showAmount")))
	s.CheckURL = params.Encode()

	tmp, _ := qrcode.Encode(fmt.Sprintf("%s?amount=%s", account.Address, s.Amount), qrcode.Low, 320)
	s.QRB64 = base64.StdEncoding.EncodeToString(tmp)

	app.render(w, http.StatusOK, "pay.html", app.newTemplateData(r))
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
