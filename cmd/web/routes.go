package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/vulkan0n/superbchat/ui"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	fileserver := http.FileServer(http.FS(ui.Files))
	r.Handle("/static/*", fileserver)

	r.Get("/", app.indexHandler)
	r.Get("/pay", app.paymentHandler)
	r.Get("/check", app.checkHandler)
	r.Get("/view", app.viewHandler)
	r.Get("/user/signup", notImplementedHandler())
	r.Post("/user/signup", notImplementedHandler())
	r.Get("/user/update", notImplementedHandler())
	r.Post("/user/update", notImplementedHandler())
	r.Get("/alert/:user/:pass", notImplementedHandler())
	r.Get("/{user}", notImplementedHandler())
	r.Post("/{user}", notImplementedHandler())

	return app.recoverPanic(app.logRequest(secureHeaders(r)))
}

func notImplementedHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Not implemented: " + r.URL.Path))
	}
}
