package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"sso/internal/grpc/auth"
	"sso/internal/services/auth"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates new gRPC server app
func New(log *slog.Logger, authService *auth.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if error occurs
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", "addr", l.Addr().String())

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
