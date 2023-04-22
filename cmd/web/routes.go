package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/justinas/alice"
	"github.com/vulkan0n/superbchat/ui"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.NotFound(app.notFound)

	fileserver := http.FileServer(http.FS(ui.Files))
	r.Handle("/static/*", fileserver)

	dynamic := alice.New(app.sessionManager.LoadAndSave)
	protected := dynamic.Append(app.requireAuthentication)

	r.Get("/", dynamic.ThenFunc(app.index).ServeHTTP)
	r.Get("/pay", dynamic.ThenFunc(app.pay).ServeHTTP)
	r.Post("/pay", dynamic.ThenFunc(app.payPost).ServeHTTP)
	r.Get("/check", dynamic.ThenFunc(app.checkHandler).ServeHTTP)
	r.Get("/view", protected.ThenFunc(app.viewHandler).ServeHTTP)
	r.Get("/user/login", dynamic.ThenFunc(app.userLogin).ServeHTTP)
	r.Post("/user/login", dynamic.ThenFunc(app.userLoginPost).ServeHTTP)
	r.Post("/user/logout", protected.ThenFunc(app.userLogoutPost).ServeHTTP)
	r.Get("/user/signup", dynamic.ThenFunc(app.userSignup).ServeHTTP)
	r.Post("/user/signup", dynamic.ThenFunc(app.userSignupPost).ServeHTTP)
	r.Get("/user/update", protected.ThenFunc(notImplementedHandler()).ServeHTTP)
	r.Post("/user/update", protected.ThenFunc(notImplementedHandler()).ServeHTTP)
	r.Get("/alert/:user/:pass", notImplementedHandler())
	r.Get("/{user}", dynamic.ThenFunc(app.superbchat).ServeHTTP)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(r)
}

func notImplementedHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Not implemented: " + r.URL.Path))
	}
}
