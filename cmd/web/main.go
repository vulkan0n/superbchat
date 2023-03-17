package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"
)

// Email settings
var enableEmail = false
var smtpHost = "smtp.purelymail.com"
var smtpPort = "587"
var smtpUser = "example@purelymail.com"
var smtpPass = "[y7EQ(xgTW_~{CUpPhO6(#"
var sendTo = []string{"example@purelymail.com"} // Comma separated recipient list

var indexTemplate *template.Template
var payTemplate *template.Template
var checkTemplate *template.Template
var alertTemplate *template.Template
var viewTemplate *template.Template
var topWidgetTemplate *template.Template
var superbchatTemplate *template.Template

type configJson struct {
	BCHAddress       string   `json:"BCHAddress"`
	MinimumDonation  float64  `json:"MinimumDonation"`
	MaxMessageChars  int      `json:"MaxMessageChars"`
	MaxNameChars     int      `json:"MaxNameChars"`
	WebViewUsername  string   `json:"WebViewUsername"`
	WebViewPassword  string   `json:"WebViewPassword"`
	OBSWidgetRefresh string   `json:"OBSWidgetRefresh"`
	Checked          bool     `json:"ShowAmountCheckedByDefault"`
	EnableEmail      bool     `json:"EnableEmail"`
	SMTPServer       string   `json:"SMTPServer"`
	SMTPPort         string   `json:"SMTPPort"`
	SMTPUser         string   `json:"SMTPUser"`
	SMTPPass         string   `json:"SMTPPass"`
	SendToEmail      []string `json:"SendToEmail"`
}

func main() {

	var conf configJson
	err := json.Unmarshal(configBytes, &conf)
	if err != nil {
		panic(err) // Fatal error, stop program
	}

	BCHAddress = conf.BCHAddress
	ScamThreshold = conf.MinimumDonation
	MessageMaxChar = conf.MaxMessageChars
	NameMaxChar = conf.MaxNameChars
	username = conf.WebViewUsername
	password = conf.WebViewPassword
	AlertWidgetRefreshInterval = conf.OBSWidgetRefresh
	enableEmail = conf.EnableEmail
	smtpHost = conf.SMTPServer
	smtpPort = conf.SMTPPort
	smtpUser = conf.SMTPUser
	smtpPass = conf.SMTPPass
	sendTo = conf.SendToEmail
	if conf.Checked == true {
		checked = " checked"
	}

	fmt.Println(BCHAddress)

	fmt.Println(fmt.Sprintf("email notifications enabled?: %t", enableEmail))
	fmt.Println(fmt.Sprintf("OBS Alert path: /alert?auth=%s", password))

	mux := http.NewServeMux()
	var styleFS = http.FS(styleFiles)
	fs := http.FileServer(styleFS)
	mux.Handle("/style/", fs)

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/superbchat", superbchatHandler)
	mux.HandleFunc("/pay", paymentHandler)
	mux.HandleFunc("/create", createHandler)
	mux.HandleFunc("/check", checkHandler)
	mux.HandleFunc("/alert", alertHandler)
	mux.HandleFunc("/view", viewHandler)
	mux.HandleFunc("/top", topwidgetHandler)

	// Create files and directory if they don't exist
	path := "log"
	_ = os.Mkdir(path, os.ModePerm)

	_, err = os.OpenFile("log/paid.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	_, err = os.OpenFile("log/alertqueue.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	_, err = os.OpenFile("log/superchats.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	indexTemplate, _ = template.ParseFS(resources, "web/index.html")
	payTemplate, _ = template.ParseFS(resources, "web/pay.html")
	checkTemplate, _ = template.ParseFS(resources, "web/check.html")
	alertTemplate, _ = template.ParseFS(resources, "web/alert.html")
	viewTemplate, _ = template.ParseFS(resources, "web/view.html")
	topWidgetTemplate, _ = template.ParseFS(resources, "web/top.html")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8900"
	}
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}

}
func mail(name string, amount string, message string) {
	body := []byte(fmt.Sprintf("From: %s\n"+
		"Subject: %s sent %s BCH\nDate: %s\n\n"+
		"%s", smtpUser, name, amount, fmt.Sprint(time.Now().Format(time.RFC1123Z)), message))

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, sendTo, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("email sent")
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
