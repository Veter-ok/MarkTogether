package app

import (
	"context"
	"net/http"

	"github.com/Veter-ok/MarkTogether/internal/api"
	"github.com/Veter-ok/MarkTogether/internal/wsserver"
)

const (
	templateDir = "../web/templates/html"
)

type App struct {
	mux    *http.ServeMux
	server *http.Server
	wsSrv  *wsserver.WSServer
	api    *api.Handler
}

func NewApp(addr string) *App {
	mux := http.NewServeMux()
	wsServer := wsserver.NewWsServer(addr)
	handler := api.NewHandler(wsServer)

	app := &App{
		mux: mux,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		wsSrv: wsServer,
		api:   handler,
	}
	app.registerRoutes()
	return app
}

func (app *App) registerRoutes() {
	app.mux.Handle("/", http.FileServer(http.Dir(templateDir)))
	app.mux.HandleFunc("/ws", app.wsSrv.WSHandler)
}

func (app *App) Start() error {
	return app.server.ListenAndServe()
}

func (a *App) Stop() error {
	a.wsSrv.Stop()
	return a.server.Shutdown(context.Background())
}
