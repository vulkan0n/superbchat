package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	embedfiles "github.com/vulkan0n/superbchat"
)

func main() {

	conf := getDefaultConfig()
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

	flag.StringVar(&BCHAddress, "addr", "bitcoincash:address", "Bitcoin Cash address to recieve founds")
	flag.Parse()

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
	logDirectory := "./cmd/log"
	_ = os.Mkdir(logDirectory, os.ModePerm)

	_, err := os.OpenFile(logDirectory+"/paid.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8900"
	}
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
}
