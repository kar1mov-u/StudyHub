package http

import (
	"StudyHub/backend/internal/modules"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/cors"

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

	srv.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://0.0.0.0:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // maximum age for preflight request cache
	}))
	srv.router.Route("/api/v1", func(r chi.Router) {
		// Module routes
		r.Get("/modules", srv.ListModulesHandler)
		r.Post("/modules", srv.CreateModuleHandler)
		r.Get("/modules/{id}", srv.GetModuleFullHandler)
		r.Delete("/modules/{id}", srv.DeleteModuleHandler)

		// Module Run routes (nested under modules)
		r.Get("/modules/{moduleID}/runs", srv.ListModuleRunsHandler)
		r.Post("/modules/{moduleID}/runs", srv.CreateModuleRunHandler)

		// Module Run routes (direct access)
		r.Get("/module-runs/{id}", srv.GetModuleRunHandler)
		r.Delete("/module-runs/{id}", srv.DeleteModuleRunHandler)

		// Academic Calendar routes
		r.Get("/academic-terms", srv.ListAcademicTermsHandler)
		r.Get("/academic-terms/active", srv.GetActiveAcademicTermHandler)
		r.Post("/academic-terms", srv.CreateAcademicTermHandler)
		r.Patch("/academic-terms/{id}/deactivate", srv.DeactivateAcademicTermHandler)
		r.Patch("/academic-terms/{id}/activate", srv.ActivateAcademicTermHandler)
	})
}

func (srv *HTTPServer) Start() {
	srv.httpServer.ListenAndServe()
}

func (srv *HTTPServer) ShutDown(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}

type Response struct {
	Data any `json:"data"`
}

type ErrResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ResponseWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	resp := Response{Data: payload}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "cannot decode JSON response", http.StatusInternalServerError)
		slog.Error("cannot encode JSON response: %s", err.Error(), err)
		return
	}
}

func ResponseWithErr(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(map[string]ErrResponse{"error": {Code: statusCode, Message: message}})
	if err != nil {
		http.Error(w, "cannot decode JSON response", http.StatusInternalServerError)
		slog.Error("cannot encode JSON response: %s", err.Error(), err)
		return
	}
}
