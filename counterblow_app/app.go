package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) StartBalancer(name string) string {
	println("Started do balancer")

	connect()
	addHit()

	println("Started proxy")

	startProxy(name)

	return fmt.Sprintf("Hello %s, started proxy!", name)
}

func (a *App) StopBalancer(name string) string {
	println("Stopped balancer")
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) AppLogAppend(name string) string {
	println("Appending...")
	return ""
}
