package connector

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"ll_test/app/pkg/httpErrors"

	llDelivery "ll_test/ll/delivery"

	"github.com/go-chi/httprate"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	err := h(w, r)
	if err != nil {
		switch e := err.(type) {
		case httpErrors.Error:
			w.WriteHeader(e.Status())
			render.JSON(w, r, httpErrors.NewRestError(e.Status(), e.Error(), e.Causes()))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, httpErrors.NewRestError(http.StatusInternalServerError, e.Error(), httpErrors.ErrInternalServerError))
		}
	}
}

// - TODO: Need to add CORS and appropriate rate limiting
func SetupRouter() *chi.Mux {
	r := chi.NewRouter()
	//- General middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(httprate.Limit(
		10,
		1*time.Minute,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			httpErrors.NewRestError(http.StatusTooManyRequests, "Too many requests", http.StatusTooManyRequests)
		}),
	)) //- 100 request per 1 minute

	// Routing
	r.Route("/littlelives", llHandler)

	return r
}

func llHandler(r chi.Router) {
	//- authentication here
	r.Group(func(r chi.Router) {
		r.Method("POST", "/user", Handler(llDelivery.LLAddNewUser))
		r.Method("POST", "/file", Handler(llDelivery.LLVideoUploadFile))
	})
}
