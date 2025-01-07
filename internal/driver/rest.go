package driver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ghazlabs/idn-remote-entry/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"gopkg.in/validator.v2"
)

type API struct {
	APIConfig
}

type APIConfig struct {
	Service core.Service `validate:"nonnil"`
}

func NewAPI(cfg APIConfig) (*API, error) {
	err := validator.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid API config: %w", err)
	}
	return &API{APIConfig: cfg}, nil
}

func (a *API) GetHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Post("/vacancies", a.serveSubmitVacancy)

	return r
}

func (a *API) serveSubmitVacancy(w http.ResponseWriter, r *http.Request) {
	var req core.SubmitRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		render.Render(w, r, NewErrorResp(NewBadRequestError(err.Error())))
		return
	}
	err = a.Service.Handle(r.Context(), req)
	if err != nil {
		render.Render(w, r, NewErrorResp(err))
		return
	}
	render.Render(w, r, NewSuccessResp(nil))
}
