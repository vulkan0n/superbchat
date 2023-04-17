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

	"github.com/skip2/go-qrcode"
	"github.com/vulkan0n/superbchat/internal/models"
	"github.com/vulkan0n/superbchat/internal/validator"
)

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

type superbchatDisplay struct {
	User    string
	MaxChar int
	MinAmnt float64
	Checked string
}

type viewPageData struct {
	ID      []string
	Name    []string
	Message []string
	Amount  []string
	Display []string
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "index.html", nil)
}

func (app *application) viewHandler(w http.ResponseWriter, r *http.Request) {
	var a viewPageData
	var displayTemp string

	u, p, ok := r.BasicAuth()
	if !ok {
		w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if (u == username) && (p == password) {
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

		csvLines, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			fmt.Println(err)
		}

		for _, line := range csvLines {
			a.ID = append(a.ID, line[0])
			a.Name = append(a.Name, line[1])
			a.Message = append(a.Message, line[2])
			a.Amount = append(a.Amount, line[3])
			displayTemp = fmt.Sprintf("<h3><b>%s</b> sent <b>%s</b> BCH:</h3><p>%s</p>", html.EscapeString(line[1]), html.EscapeString(line[3]), line[2])
			a.Display = append(a.Display, displayTemp)
		}

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return // return http 401 unauthorized error
	}
	reverse(a.Display)

	files := []string{
		"./ui/html/pages/view.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = ts.Execute(w, a)
	if err != nil {
		app.errorLog.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) checkHandler(w http.ResponseWriter, r *http.Request) {

	var c checkPage
	c.Meta = `<meta http-equiv="Refresh" content="5">`
	c.Addy = BCHAddress
	c.Received, _ = strconv.ParseFloat(r.FormValue("amount"), 64)
	c.Name = truncateStrings(r.FormValue("name"), NameMaxChar)
	c.Msg = truncateStrings(r.FormValue("msg"), MessageMaxChar)
	c.Receipt = "Waiting for payment..."

	var txsWallet []string
	getTXs(&txsWallet)
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
					if c.Received >= ScamThreshold {
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
	if received < ScamThreshold {
		*receiptPtr = "<b style='color:red'>Scammed! " + fmt.Sprint(received) + " is below minimum</b>"
	} else {
		*receiptPtr = "<b>" + fmt.Sprint(received) + " BCH Received! Superchat sent</b>"
	}
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
	if (u == username) && (p == password) {
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

func (app *application) paymentHandler(w http.ResponseWriter, r *http.Request) {
	if BCHAddress != "" {
		var s superChat
		s.Amount = html.EscapeString(r.FormValue("amount"))
		if r.FormValue("amount") == "" {
			s.Amount = fmt.Sprint(ScamThreshold)
		}
		if r.FormValue("name") == "" {
			s.Name = "Anonymous"
		} else {
			s.Name = html.EscapeString(truncateStrings(condenseSpaces(r.FormValue("name")), NameMaxChar))
		}
		s.Message = html.EscapeString(truncateStrings(condenseSpaces(r.FormValue("message")), MessageMaxChar))
		s.Address = BCHAddress

		params := url.Values{}
		params.Add("amount", s.Amount)
		params.Add("name", s.Name)
		params.Add("msg", r.FormValue("message"))
		params.Add("show", html.EscapeString(r.FormValue("showAmount")))
		s.CheckURL = params.Encode()

		tmp, _ := qrcode.Encode(fmt.Sprintf("%s?amount=%s", BCHAddress, s.Amount), qrcode.Low, 320)
		s.QRB64 = base64.StdEncoding.EncodeToString(tmp)

		files := []string{
			"./ui/html/base.html",
			"./ui/html/pages/pay.html",
		}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.errorLog.Fatal(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		err = ts.ExecuteTemplate(w, "base", s)
		if err != nil {
			app.errorLog.Fatal(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return // return http 401 unauthorized error
	}
}

type userSignupForm struct {
	User                string `form:"user"`
	Password            string `form:"password"`
	RepeatedPassword    string `form:"repeated-password"`
	Address             string `form:"address"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {

	form := userSignupForm{
		User:             "",
		Password:         "",
		RepeatedPassword: "",
		Address:          "",
	}

	data := app.newTemplateData(r)
	data.Form = form
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
			return
		} else {
			app.serverError(w, err)
		}
	}
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, fmt.Sprintf("/%s", form.User), http.StatusSeeOther)
}

func (app *application) superbchatHandler(w http.ResponseWriter, r *http.Request) {
	var user = "vulkan0n"
	var s superbchatDisplay
	s.User = user
	s.MaxChar = MessageMaxChar
	s.MinAmnt = ScamThreshold
	s.Checked = checked

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/header.html",
		"./ui/html/pages/superbchat.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", s)
	if err != nil {
		app.errorLog.Fatal(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
