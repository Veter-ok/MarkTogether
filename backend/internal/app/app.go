package app

import (
	"context"
	"net/http"

	"github.com/Veter-ok/MarkTogether/internal/api"
	"github.com/Veter-ok/MarkTogether/internal/document"
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
	store  document.Store
}

func NewApp(addr string) *App {
	mux := http.NewServeMux()
	store := document.NewInMemoryStore()
	wsServer := wsserver.NewWsServer(addr, store)
	handler := api.NewHandler(wsServer)

	app := &App{
		mux: mux,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		wsSrv: wsServer,
		api:   handler,
		store: store,
	}
	app.registerRoutes()
	return app
}

func (app *App) registerRoutes() {
	app.mux.Handle("/", http.FileServer(http.Dir(templateDir)))

	docHandler := api.NewDocumentHandler(app.store)
	app.mux.HandleFunc("POST /api/documents", docHandler.CreateDocument)
	app.mux.HandleFunc("GET /api/documents/{id}", docHandler.GetDocument)

	app.mux.HandleFunc("/ws", app.wsSrv.WSHandler)
}

func (app *App) Start() error {
	return app.server.ListenAndServe()
}

func (a *App) Stop() error {
	a.wsSrv.Stop()
	return a.server.Shutdown(context.Background())
}
