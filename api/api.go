package api

import (
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/audit"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"net/http"
	"runtime/debug"
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
	app             interface{}
	authService     string
	premissions     interface{}
	singleUserToken string
	SolidAuth       bool
	logger          *mlog.Logger
	audit           *audit.Audit
}

func NewAPI(app interface{}, singleUserToken string, authService string, permissions interface{},
	logger *mlog.Logger, audit *audit.Audit) *API {
	return &API{
		app:             app,
		authService:     authService,
		premissions:     permissions,
		singleUserToken: singleUserToken,
		logger:          logger,
		audit:           audit,
	}
}

func (a *API) RegisterRoutes(r *mux.Router) {
	apiv1 := r.PathPrefix("/api/v1").Subrouter()
	apiv1.Use(a.panicHandler)
	apiv1.Use(a.requireCSRFToken)

	// Auth APIs
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
				// a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *API) requireCSRFToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.checkCSRFToken(r) {
			a.logger.Error("checkCSRFToken FAILED")
			//a.errorResponse(w, r.URL.Path, http.StatusBadRequest, "checkCSRFToken FAILED", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
