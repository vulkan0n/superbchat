package main

import (
	"net/http"

	embedfiles "github.com/vulkan0n/superbchat"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	var styleFS = http.FS(embedfiles.StyleFiles)
	fs := http.FileServer(styleFS)
	mux.Handle("/ui/static/", fs)

	mux.HandleFunc("/", app.indexHandler)
	mux.HandleFunc("/superbchat", app.superbchatHandler)
	mux.HandleFunc("/pay", app.paymentHandler)
	mux.HandleFunc("/create", app.createHandler)
	mux.HandleFunc("/check", app.checkHandler)
	mux.HandleFunc("/alert", app.alertHandler)
	mux.HandleFunc("/view", app.viewHandler)
	mux.HandleFunc("/top", app.topwidgetHandler)

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
