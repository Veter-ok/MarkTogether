package api

import "github.com/Veter-ok/MarkTogether/internal/wsserver"

type Handler struct {
	hub *wsserver.WSServer
}

func NewHandler(wsSrv *wsserver.WSServer) *Handler {
	return &Handler{wsSrv}
}
