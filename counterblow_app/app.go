package main

import (
	"context"
	"errors"
	"fmt"
	"time"

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
}

func (a *App) domready(ctx context.Context) {
	println("DOM ready")
	a.RefreshRules()
}

func (a *App) OnDOMContentLoaded(arg1 string) {
	//a.RefreshRules()
}

func (a *App) RefreshRules() {
	time.Sleep(1 * time.Second) // there is a bug: DOM is not actually ready
	TextAreaLog("Application started!")
	TextAreaLog("Loading rules...")

	err := database_connect()
	if err != nil {
		TextAreaLog(err.Error())
	}

	runtime.EventsEmit(global_app.ctx, "rcv:clear_rules_listbox")
	rules, err = database_loadRules()

	for rule_id, rule := range rules {
		fmt.Printf("Loaded rule: %v\n", rule)
		runtime.EventsEmit(global_app.ctx, "rcv:add_served_rule", rule.rule_id, rule.rule_type, rule.rule_ipaddr, rule.rule_subnetmask, rule.rule_servers, rule.rule_source, rule.rule_dest)
		TextAreaLog(fmt.Sprintf("Loaded rule %d", rule_id))
	}
	TextAreaLog("Loading rules ended.")
}

func (a *App) StartBalancer(bindIp string, port int) error {
	// database_addHit("from", "to")
	if len(rules) == 0 {
		TextAreaLog("No routing rules!")
		return errors.New("No routing rules!")
	}
	TextAreaLog("Starting load balancer...")
	err = startReverseProxy(bindIp, port, rules)
	TextAreaLog("Error: " + err.Error())
	return nil
}

func (a *App) StopBalancer(name string) string {
	println("Stopped balancer")
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) AddRule(rule_type int, ip string, mask int, servers string, rule_source string, rule_dest string) {
	// saving rule to db. Called from frontend
	// prepare rule for insertion
	println("Adding a rule...")
	var rule RoutingRule
	rule.rule_type = rule_type
	rule.rule_ipaddr = ip
	rule.rule_subnetmask = mask
	rule.rule_servers = servers
	rule.rule_source = rule_source
	rule.rule_dest = rule_dest
	database_addRule(rule)
	println("Added")

}

func (a *App) RemoveRule(rule_id int) {
	// saving rule to db. Called from frontend
	// prepare rule for insertion
	println("Removing a rule...")
	err = database_removeRule(rule_id)
	if err != nil {
		TextAreaLog("Correctly removed item")
	} else {
		TextAreaLog(err.Error())
	}
	println("Removed.")

}

func UpdateServedPages(count int) {
	runtime.EventsEmit(global_app.ctx, "rcv:update_served_pages", count)

}

func TextAreaLog(appendText string) {
	runtime.EventsEmit(global_app.ctx, "rcv:add_log_string", appendText+"\n")
}
