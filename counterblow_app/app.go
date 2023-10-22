package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct

type App struct {
	ctx context.Context
}

var global_app *App

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	global_app = a
}

func (a *App) StartBalancer(name string) string {
	println("Started do balancer")

	database_connect()
	database_addHit("from", "to")

	println("Started proxy")

	//startHttpServer(name) // era solo per debug
	startReverseProxy("0.0.0.0", 8080)

	return fmt.Sprintf("Hello %s, started proxy!", name)

}

func (a *App) StopBalancer(name string) string {
	println("Stopped balancer")
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) AppLogAppend(name string) {
	println("Appending...")
}

func UpdateServedPages(count int) {
	//go func() {
	runtime.EventsEmit(global_app.ctx, "rcv:update_served_pages", count)

	//}()
}

func TextAreaLog(appendText string) {
	//go func() {
	runtime.EventsEmit(global_app.ctx, "rcv:add_log_string", appendText+"\n")

	//}()
}
