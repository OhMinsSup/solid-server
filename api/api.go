package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/shared/mlog"
	"net/http"
	"runtime/debug"
	"solid-server/app"
	"solid-server/model"
	"solid-server/services/auth"
	"solid-server/utils"
)

const (
	HeaderRequestedWith    = "X-Requested-With"
	HeaderRequestedWithXML = "XMLHttpRequest"
	HeaderRequestedWithKy  = "ky"
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
	app         *app.App
	authService string
	premissions interface{}
	logger      *mlog.Logger
}

func NewAPI(app *app.App, authService string, permissions interface{},
	logger *mlog.Logger) *API {
	return &API{
		app:         app,
		authService: authService,
		premissions: permissions,
		logger:      logger,
	}
}

func (a *API) RegisterRoutes(r *mux.Router) {
	apiv1 := r.PathPrefix("/api/v1").Subrouter()
	apiv1.Use(a.panicHandler)
	apiv1.Use(a.requireCSRFToken)

	// Auth APIs
	apiv1.HandleFunc("/auth/login", a.handleLogin).Methods("POST")
	apiv1.HandleFunc("/auth/register", a.handleRegister).Methods("POST")

	// User APIs
	apiv1.HandleFunc("/users/me", a.requiredAuth(a.handleUserMe)).Methods("GET")

	// Post APIs
	apiv1.HandleFunc("/posts", a.requiredAuth(a.handleCreatePost)).Methods("POST")
	apiv1.HandleFunc("/posts/{postID}", a.handleReadPost).Methods("GET")

}

func (a *API) checkCSRFToken(r *http.Request) bool {
	token := r.Header.Get(HeaderRequestedWith)
	return token == HeaderRequestedWithXML || token == HeaderRequestedWithKy
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

func (a *API) requiredAuth(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return a.authMiddleware(handler, true)
}

func (a *API) authMiddleware(handler func(w http.ResponseWriter, r *http.Request), required bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, _ := auth.ParseAuthTokenFromRequest(r)
		if len(token) == 0 {
			a.errorResponse(w, r.URL.Path, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		a.logger.Debug(`attachAuthMiddleware`, mlog.Bool("accessToken", len(token) > 0))
		decodeTokenData, _ := a.app.GetAuth().VerifyAccessToken(token)
		if decodeTokenData == nil {
			if required {
				a.errorResponse(w, r.URL.Path, http.StatusUnauthorized, "Unauthorized", nil)
				return
			}
			handler(w, r)
			return
		}

		user, err := a.app.GetUser(decodeTokenData.UserID)
		if err != nil {
			if required {
				a.errorResponse(w, r.URL.Path, http.StatusUnauthorized, "", err)
				return
			}
			handler(w, r)
			return
		}

		authService := user.AuthService
		if authService != a.authService {
			a.logger.Error(`Session authService mismatch`,
				mlog.String("userID", user.ID),
				mlog.String("want", a.authService),
				mlog.String("got", authService),
			)
			a.errorResponse(w, r.URL.Path, http.StatusUnauthorized, "", err)
			return
		}

		ctx := context.WithValue(r.Context(), sessionContextKey, user)
		handler(w, r.WithContext(ctx))
	}
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
