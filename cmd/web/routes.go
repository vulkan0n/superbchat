package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/justinas/alice"
	"github.com/vulkan0n/superbchat/ui"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	fileserver := http.FileServer(http.FS(ui.Files))
	r.Handle("/static/*", fileserver)

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	r.Get("/", dynamic.ThenFunc(app.index).ServeHTTP)
	r.Get("/pay", dynamic.ThenFunc(app.paymentHandler).ServeHTTP)
	r.Get("/check", dynamic.ThenFunc(app.checkHandler).ServeHTTP)
	r.Get("/view", dynamic.ThenFunc(app.viewHandler).ServeHTTP)
	r.Get("/user/login", dynamic.ThenFunc(app.userLogin).ServeHTTP)
	r.Post("/user/login", dynamic.ThenFunc(app.userLoginPost).ServeHTTP)
	r.Post("/user/logout", dynamic.ThenFunc(app.userLogoutPost).ServeHTTP)
	r.Get("/user/signup", dynamic.ThenFunc(app.userSignup).ServeHTTP)
	r.Post("/user/signup", dynamic.ThenFunc(app.userSignupPost).ServeHTTP)
	r.Get("/user/update", notImplementedHandler())
	r.Post("/user/update", notImplementedHandler())
	r.Get("/alert/:user/:pass", notImplementedHandler())
	r.Get("/{user}", notImplementedHandler())
	r.Post("/{user}", notImplementedHandler())

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(r)
}

func notImplementedHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Not implemented: " + r.URL.Path))
	}
}
