package main

import (
	"net/http"

	"github.com/labstack/echo/v4/middleware"
	"github.com/vulkan0n/superbchat/ui"
)

func (app *application) routes() {
	// r := chi.NewRouter()
	// r.NotFound(app.notFound)

	// fileserver := http.FileServer(http.FS(ui.Files))
	// r.Handle("/static/*", fileserver)

	// dynamic := alice.New(app.sessionManager.LoadAndSave)
	// protected := dynamic.Append(app.requireAuthentication)

	// r.Get("/", dynamic.ThenFunc(app.index).ServeHTTP)
	// r.Post("/pay", dynamic.ThenFunc(app.payPost).ServeHTTP)
	// r.Get("/check/{superchatId:[0-9]+}/{accountId:[0-9]+}", dynamic.ThenFunc(app.check).ServeHTTP)
	// r.Get("/view", protected.ThenFunc(app.view).ServeHTTP)
	// r.Get("/user/login", dynamic.ThenFunc(app.userLogin).ServeHTTP)
	// r.Post("/user/login", dynamic.ThenFunc(app.userLoginPost).ServeHTTP)
	// r.Post("/user/logout", protected.ThenFunc(app.userLogoutPost).ServeHTTP)
	// r.Get("/user/signup", dynamic.ThenFunc(app.userSignup).ServeHTTP)
	// r.Post("/user/signup", dynamic.ThenFunc(app.userSignupPost).ServeHTTP)
	// r.Get("/user/settings", protected.ThenFunc(app.settings).ServeHTTP)
	// r.Post("/user/settings", protected.ThenFunc(app.settingsPost).ServeHTTP)
	// r.Get("/alert/{token}", dynamic.ThenFunc(app.alert).ServeHTTP)
	// r.Get("/{user}", dynamic.ThenFunc(app.superbchat).ServeHTTP)

	// standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// return standard.Then(r)

	app.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}\n",
	}))

	app.echo.POST("/user-signup", app.postUserSignup)
	app.echo.POST("/user-login", app.postUserLogin)

	app.echo.POST("/test", app.postTest)

	app.echo.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "frontend/dist",
		Index:      "index.html", // This is the default html page
		Browse:     false,
		HTML5:      true,
		Filesystem: http.FS(ui.Frontend),
	}))
}
