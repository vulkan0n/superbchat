package main

import (
	"net/http"

	"github.com/labstack/echo/v4/middleware"
	"github.com/vulkan0n/superbchat/ui"
)

func (app *application) routes() {
	app.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}\n",
	}))

	app.echo.POST("/user-signup", app.postUserSignup)
	app.echo.POST("/user-login", app.postUserLogin)
	app.echo.POST("/verify-tkn", app.postVerifyToken)
	app.echo.GET("/user-id/:id", app.getUserInfo)
	app.echo.GET("/user/:user", app.getUserInfoByName)
	app.echo.POST("/superbchat", app.postSuperbchat)
	app.echo.DELETE("/superbchat/:id", app.deleteSuperchat)
	app.echo.POST("/superbchat-get", app.postSuperchatsByTkn)
	app.echo.POST("/settings", app.postUserSettingsByTkn)

	app.echo.GET("/ws", app.wsHub.handleWebSocket)

	app.echo.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "frontend/dist",
		Index:      "index.html",
		Browse:     false,
		HTML5:      true,
		Filesystem: http.FS(ui.Frontend),
	}))
}
