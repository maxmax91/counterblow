package main

import (
	"context"
	"fmt"

	"github.com/labstack/gommon/log"
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
	TextAreaLog("Application started")
	//database_loadRules()

}

func (a *App) OnDOMContentLoaded(arg1 string) {
	println("DOM loaded")
	TextAreaLog("Application started!")
	TextAreaLog("Loading rules...")

	err := database_connect()
	if err != nil {
		TextAreaLog(err.Error())
	}

	rules := database_loadRules()
	for _, rule := range rules {
		log.Info(fmt.Printf("%v", rule))
	}
}

func (a *App) StartBalancer(bindIp string, port int) bool {
	println("Started do balancer")
	TextAreaLog("Loading rules...")

	// database_addHit("from", "to")

	//startHttpServer(name) // era solo per debug
	TextAreaLog("Starting load balancer...")
	startReverseProxy(bindIp, port)
	return true
}

func (a *App) StopBalancer(name string) string {
	println("Stopped balancer")
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SaveRule(rule_type int, ip string, mask int, servers string) {
	// saving rule to db
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
