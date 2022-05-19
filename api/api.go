package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"net/http"
	"runtime/debug"
	"solid-server/app"
	"solid-server/model"
	"solid-server/utils"
)

const (
	HeaderRequestedWith    = "X-Requested-With"
	HeaderRequestedWithXML = "XMLHttpRequest"
	UploadFormFileKey      = "file"
)

const (
	ErrorNoTeamCode    = 1000
	ErrorNoTeamMessage = "No team"
)

type Permission struct {
	msg string
}

func (p Permission) Error() string {
	return p.msg
}

// REST APIs

type API struct {
	app             *app.App
	authService     string
	premissions     interface{}
	SolidAuth       bool
	logger          *mlog.Logger
}

func NewAPI(app *app.App, authService string, permissions interface{},
	logger *mlog.Logger) *API {
	return &API{
		app:             app,
		authService:     authService,
		premissions:     permissions,
		logger:          logger,
	}
}

func (a *API) RegisterRoutes(r *mux.Router) {
	apiv1 := r.PathPrefix("/api/v1").Subrouter()
	apiv1.Use(a.panicHandler)
	apiv1.Use(a.requireCSRFToken)

	// Auth APIs
	apiv1.HandleFunc("/auth/login", a.handleLogin).Methods("POST")
	apiv1.HandleFunc("/auth/register", a.handleRegister).Methods("POST")
}

func (a *API) checkCSRFToken(r *http.Request) bool {
	token := r.Header.Get(HeaderRequestedWith)
	return token == HeaderRequestedWithXML
}

func (a *API) panicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); !utils.IsNilFixed(err) {
				a.logger.Error("Http handler panic",
					mlog.Any("panic", err),
					mlog.String("stack", string(debug.Stack())),
					mlog.String("uri", r.URL.Path),
				)
				a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *API) requireCSRFToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.checkCSRFToken(r) {
			a.logger.Error("checkCSRFToken FAILED")
			a.errorResponse(w, r.URL.Path, http.StatusBadRequest, "checkCSRFToken FAILED", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Response helpers

func (a *API) errorResponse(w http.ResponseWriter, api string, code int, message string, sourceError error) {
	if code == http.StatusUnauthorized || code == http.StatusForbidden {
		a.logger.Debug("API DEBUG",
			mlog.Int("code", code),
			mlog.Err(sourceError),
			mlog.String("msg", message),
			mlog.String("api", api),
		)
	} else {
		a.logger.Error("API ERROR",
			mlog.Int("code", code),
			mlog.Err(sourceError),
			mlog.String("msg", message),
			mlog.String("api", api),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(model.ErrorResponse{Error: message, ErrorCode: code})
	if err != nil {
		data = []byte("{}")
	}
	w.WriteHeader(code)
	_, _ = w.Write(data)
}

func jsonStringResponse(w http.ResponseWriter, code int, message string) { //nolint:unparam
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

func jsonBytesResponse(w http.ResponseWriter, code int, json []byte) { //nolint:unparam
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(json)
}
