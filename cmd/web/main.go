package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	embedfiles "github.com/vulkan0n/superbchat"
)

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
	err := json.Unmarshal(embedfiles.ConfigBytes, &conf)
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
	var styleFS = http.FS(embedfiles.StyleFiles)
	fs := http.FileServer(styleFS)
	mux.Handle("/ui/static/", fs)

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/superbchat", superbchatHandler)
	mux.HandleFunc("/pay", paymentHandler)
	mux.HandleFunc("/create", createHandler)
	mux.HandleFunc("/check", checkHandler)
	mux.HandleFunc("/alert", alertHandler)
	mux.HandleFunc("/view", viewHandler)
	mux.HandleFunc("/top", topwidgetHandler)

	// Create files and directory if they don't exist
	logDirectory := "../../log"
	_ = os.Mkdir(logDirectory, os.ModePerm)

	_, err = os.OpenFile(logDirectory+"/paid.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	_, err = os.OpenFile(logDirectory+"/alertqueue.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	_, err = os.OpenFile(logDirectory+"/superchats.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	indexTemplate, _ = template.ParseFS(embedfiles.Resources, "ui/html/index.html")
	payTemplate, _ = template.ParseFS(embedfiles.Resources, "ui/html/pay.html")
	checkTemplate, _ = template.ParseFS(embedfiles.Resources, "ui/html/check.html")
	alertTemplate, _ = template.ParseFS(embedfiles.Resources, "ui/html/alert.html")
	viewTemplate, _ = template.ParseFS(embedfiles.Resources, "ui/html/view.html")
	topWidgetTemplate, _ = template.ParseFS(embedfiles.Resources, "ui/html/top.html")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8900"
	}
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}

}
