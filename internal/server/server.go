// Модуль server служит для запуска сервера с указанием http методов.
package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/nasik90/gophkeeper/internal/common/logger"
	"github.com/nasik90/gophkeeper/internal/server/handler"
	"go.uber.org/zap"

	middleware "github.com/nasik90/gophkeeper/internal/server/middlewares"
)

// Server - структура, которая характеризует сервер.
// Содержит встроенную структуру из типовой библиотеки http.Server и handler.
type Server struct {
	http.Server
	handler *handler.Handler
}

// NewServer создает экземпляр структуры Server.
func NewServer(handler *handler.Handler, serverAddress string) *Server {
	s := &Server{}
	s.Addr = serverAddress
	s.handler = handler
	return s
}

// RunServer запускает сервер.
func (s *Server) RunServer() error {

	logger.Log.Info("Running server", zap.String("address", s.Addr))

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Post("/user/register", s.handler.RegisterNewUser())
		r.Post("/user/login", s.handler.LoginUser())
		r.Post("/user/loadSecret", middleware.Auth(s.handler.LoadSecret()))
	})
	s.Handler = logger.RequestLogger((r.ServeHTTP))
	var err error
	// if s.enableHTTPS {
	// 	err = s.ListenAndServeTLS("server.crt", "server.key")
	// } else {
	err = s.ListenAndServe()
	// }
	if err != nil {
		return err
	}

	return nil
}

// StopServer останавливает сервер.
func (s *Server) StopServer(ctx context.Context) error {
	return s.Shutdown(ctx)
}
