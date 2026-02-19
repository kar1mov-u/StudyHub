package http

import (
	"StudyHub/internal/auth"
	"StudyHub/internal/modules"
	"StudyHub/internal/resources"
	"StudyHub/internal/users"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi"
)

type HTTPServer struct {
	moduleSrv   *modules.ModuleService
	authSrv     *auth.AuthService
	userSrv     *users.UserService
	resourceSrv *resources.ResourceService
	httpServer  *http.Server
	router      *chi.Mux
}

func NewHTTPServer(moduleSrv *modules.ModuleService, userSrv *users.UserService, authSrv *auth.AuthService, resSrv *resources.ResourceService, port string) *HTTPServer {
	router := chi.NewMux()
	s := HTTPServer{
		moduleSrv:   moduleSrv,
		userSrv:     userSrv,
		authSrv:     authSrv,
		resourceSrv: resSrv,
		router:      router,
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
		AllowedOrigins:   []string{"http://localhost:80", "http://0.0.0.0:80"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // maximum age for preflight request cache
	}))
	srv.router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		//auth routes
		r.Group(func(pub chi.Router) {
			pub.Get("/health", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("dev setup is not working"))
			})
			pub.Post("/auth/login", srv.LoginHandler)
			pub.Post("/users", srv.CreateUserHandler)

		})

		//User routes
		r.Group(func(priv chi.Router) {
			priv.Use(srv.authSrv.JWTMiddleware)

			priv.Get("/users", srv.ListUsersHandler)
			priv.Get("/users/me", srv.GetMeHandler)
			priv.Get("/users/{id}", srv.GetUserHandler)
			priv.Delete("/users/{id}", srv.DeleteUserHandler)
			// Module routes
			priv.Get("/modules", srv.ListModulesHandler)
			priv.Post("/modules", srv.CreateModuleHandler)
			priv.Get("/modules/{id}", srv.GetModuleFullHandler)
			priv.Delete("/modules/{id}", srv.DeleteModuleHandler)

			// Module Run routes (nested under modules)
			priv.Get("/modules/{moduleID}/runs", srv.ListModuleRunsHandler)
			priv.Post("/modules/{moduleID}/runs", srv.CreateModuleRunHandler)

			// Module Run routes (direct access)
			priv.Get("/module-runs/{id}", srv.GetModuleRunHandler)
			priv.Delete("/module-runs/{id}", srv.DeleteModuleRunHandler)

			// Academic Calendar routes
			//do not like how the academic terms are build, i should be able to just get the current one.
			priv.Get("/academic-terms/current", srv.GetActiveAcademicTermHandler)
			priv.Post("/academic-terms/new-term", srv.CreateAcademicTermHandler)

			//Resources routes
			priv.Post("/resources/file/{week_id}", srv.UploadFileHandler)
			priv.Post("/resources/link/{week_id}", srv.CreateLinkResource)
			priv.Delete("/resources/{id}", srv.DeleteResourceHandler)
			priv.Get("/resources/{id}", srv.GetResourceHandler)
			priv.Get("/resources/weeks/{week_id}", srv.ListResourcesForWeekHandler)
			priv.Get("/resources/users/{user_id}", srv.ListResourcesForUserHandler)

			//interanl processes
		})
		r.Group(func(inter chi.Router) {
			inter.Get("/interal/jobs/cleanup-orphaned-objects", srv.CleanOrphanObjectsHandler)

		})
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
