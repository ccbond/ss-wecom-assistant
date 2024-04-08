// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package server

import (
	"ss-wecom-assistant/internal/config"
	"ss-wecom-assistant/internal/logger"
	"ss-wecom-assistant/internal/services"

	"github.com/SyntSugar/ss-infra-go/api/server"
	"github.com/SyntSugar/ss-infra-go/log"
	"go.uber.org/zap"
)

type Services struct {
	WechatService services.WechatService
	ChatService   services.ChatService
}

type Server struct {
	apiServer *server.Server
	logger    *log.Logger
	config    *config.Config
	svcs      *Services
}

// NewServer creates a new server instance.
func NewServer(cfg *config.Config, services *Services) (*Server, error) {
	if err := logger.Init(cfg.LogConfig.LogLevel); err != nil {
		return nil, err
	}

	apiServer, err := server.New(&server.Config{
		API:   cfg.API,
		Admin: cfg.Admin,
	}, logger.Get())
	if err != nil {
		return nil, err
	}

	return &Server{
		apiServer: apiServer,
		config:    cfg,
		logger:    logger.Get(),
		svcs:      services,
	}, nil
}

// Run starts the server.
func (srv *Server) Run() error {
	setupAPIRouters(srv)
	if err := srv.apiServer.Run(); err != nil {
		return err
	}

	logger.Get().With(
		zap.Any("api", srv.config.API),
		zap.Any("admin", srv.config.Admin),
	).Info("The server was listening")

	return nil
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (srv *Server) Shutdown() {
	if err := srv.apiServer.Shutdown(); err != nil {
		srv.logger.With(zap.String("err", err.Error())).Error("Shutdown error")
	}
	srv.logger.Info("The server was shutdown normally, see you lala.")
}
