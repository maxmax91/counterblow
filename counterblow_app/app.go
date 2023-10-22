package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct

var domLoaded bool
var rules []RoutingRule

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

func AddLoadedRule(rule RoutingRule) {

}

func (a *App) OnDOMContentLoaded(arg1 string) {
	if domLoaded == true {
		TextAreaLog("Application started!")
		TextAreaLog("Loading rules...")

		err := database_connect()
		if err != nil {
			TextAreaLog(err.Error())
		}

		runtime.EventsEmit(global_app.ctx, "rcv:clear_rules_listbox")
		rules = database_loadRules()

		for _, rule := range rules {
			fmt.Printf("Loaded rule: %v\n", rule)
			runtime.EventsEmit(global_app.ctx, "rcv:add_served_rule", rule.rule_id, rule.rule_type, rule.rule_ipaddr, rule.rule_subnetmask, rule.rule_servers)
		}

	} else {
		fmt.Printf("Waiting second dom loading")
		domLoaded = true
	}
}

func (a *App) StartBalancer(bindIp string, port int) bool {
	println("Started do balancer")
	TextAreaLog("Loading rules...")

	// database_addHit("from", "to")

	TextAreaLog("Starting load balancer...")
	startReverseProxy(bindIp, port, rules)
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
