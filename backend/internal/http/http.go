package http

import (
	"StudyHub/backend/internal/modules"
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

type HTTPServer struct {
	moduleSrv  *modules.ModuleService
	httpServer *http.Server
	router     *chi.Mux
}

func NewHTTPServer(moduleSrv *modules.ModuleService, port string) *HTTPServer {
	router := chi.NewMux()
	s := HTTPServer{
		moduleSrv: moduleSrv,
		router:    router,
		httpServer: &http.Server{
			Addr:    port,
			Handler: router,
		},
	}
	s.registerRoutes()
	return &s
}

func (srv *HTTPServer) registerRoutes() {
	srv.router.Route("/api/v1", func(r chi.Router) {

		// r.Get("/modules", )
	})
}

func (srv *HTTPServer) Start() {
	srv.httpServer.ListenAndServe()
}

func (srv *HTTPServer) ShutDown(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}
