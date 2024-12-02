package http

import (
	"fmt"
	"net/http"
	"net/rpc"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	logfile "github.com/mkadit/go-toybox/internal/logger"
	"github.com/mkadit/go-toybox/internal/models"
	"github.com/mkadit/go-toybox/internal/ports"
)

// Adapter implements the HTTP interface
type Adapter struct {
	api    ports.APIPort
	conf   models.Configuration
	router *chi.Mux
	rpc    *rpc.Server
}

// NewAdapter creates a new Adapter
func NewAdapter(api ports.APIPort, conf models.Configuration) *Adapter {
	return &Adapter{api: api, conf: conf}
}

func (ad Adapter) Setup() (err error) {
	logfile.LogEvent(fmt.Sprintf("serving on route: %d", ad.conf.Server.Port))
	// l, err := net.Listen("tcp", fmt.Sprintf(":%d", ad.conf.Server.Port))
	// if err != nil {
	// 	return err
	// }
	// http.Serve(l, ad.router)

	err = http.ListenAndServe(fmt.Sprintf(":%d", ad.conf.Server.Port), ad.router)
	if err != nil {
		return err
	}
	return
}

// the specified port
func (ad *Adapter) SetupRouter() {
	r := chi.NewRouter()
	r.Use(
		middleware.Recoverer,
		cors.Handler(
			cors.Options{
				AllowedOrigins: []string{"https://*", "http://*"},
				// AllowOriginFunc: func(r *http.Request, origin string) bool {
				// },
				AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut},
				AllowedHeaders: []string{
					"Content-Length",
					"Content-Type",
					"Accept",
					"Authorization"},
				ExposedHeaders: []string{
					"Content-Length",
					"Content-Type",
					"Accept",
					"X-Token",
					"X-Username",
				},
				AllowCredentials: true,
				// MaxAge:             0,
				// OptionsPassthrough: false,
				// Debug: false,
			}),
		LogMiddleware,
	)

	r.Get("/", ad.Echo)
	testRoute := chi.NewRouter()

	testRoute.Route("/test", func(r chi.Router) {
		// r.Use(ValidateEcho)
		r.Get("/", ad.Echo)
	})
	r.Mount("/api", testRoute)

	// if ad.rpc != nil {
	//
	// 	ad.rpc.HandleHTTP("/rpc", "/dprc")
	// 	r.Handle("/rpc", ad.rpc)
	// 	r.Handle("/debug", ad.rpc)
	// }

	ad.router = r
	return

}

func (ad *Adapter) SetupRPC(rpcs *rpc.Server) (err error) {
	ad.rpc = rpcs
	return
}
