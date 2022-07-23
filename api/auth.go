package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"solid-server/services/auth"
	"strings"
	"time"
)

const (
	MinimumPasswordLength = 8
)

type ParamError struct {
	msg string
}

func (pe ParamError) Error() string {
	return pe.msg
}

func isValidPassword(password string) error {
	if len(password) < MinimumPasswordLength {
		return ParamError{fmt.Sprintf("password must be at least %d characters", MinimumPasswordLength)}
	}
	return nil
}

// LoginRequest 로그인 요청 Body
// swagger:model
type LoginRequest struct {
	// Type of login, currently must be set to "normal"
	// required: true
	Type string `json:"type"`

	// If specified, login using username
	// required: false
	Username string `json:"username"`

	// If specified, login using email
	// required: false
	Email string `json:"email"`

	// Password
	// required: true
	Password string `json:"password"`

	// MFA token
	// required: false
	// swagger:ignore
	MfaToken string `json:"mfa_token"`
}

// LoginResponse 로그인 Response
// swagger:model
type LoginResponse struct {
	// Session token
	// required: true
	Token string `json:"token"`
}

func LoginResponseFromJSON(data io.Reader) (*LoginResponse, error) {
	var resp LoginResponse
	if err := json.NewDecoder(data).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RegisterRequest 회원가입 요청 Body
// swagger:model
type RegisterRequest struct {
	// User name
	// required: true
	Username string `json:"username"`

	// User's email
	// required: true
	Email string `json:"email"`

	// Password
	// required: true
	Password string `json:"password"`
}

func (rd *RegisterRequest) IsValid() error {
	if strings.TrimSpace(rd.Username) == "" {
		return ParamError{"username is required"}
	}
	if strings.TrimSpace(rd.Email) == "" {
		return ParamError{"email is required"}
	}
	if !auth.IsEmailValid(rd.Email) {
		return ParamError{"invalid email format"}
	}
	if rd.Password == "" {
		return ParamError{"password is required"}
	}
	return isValidPassword(rd.Password)
}

func (a *API) handleLogin(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /login login
	//
	// Login user
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   description: Login request
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/LoginRequest"
	// responses:
	//   '200':
	//     description: success
	//     schema:
	//       "$ref": "#/definitions/LoginResponse"
	//   '401':
	//     description: invalid login
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	//   '500':
	//     description: internal error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "", err)
		return
	}

	var loginData LoginRequest
	err = json.Unmarshal(requestBody, &loginData)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "", err)
		return
	}

	if loginData.Type == "normal" {
		token, err := a.app.Login(loginData.Username, loginData.Email, loginData.Password, loginData.MfaToken)
		if err != nil {
			a.errorResponse(w, r.URL.Path, http.StatusUnauthorized, err.Error(), err)
			return
		}
		json, err := json.Marshal(LoginResponse{Token: token})
		if err != nil {
			a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "", err)
			return
		}

		// cookie will get expired after 1 days
		expires := time.Now().AddDate(0, 0, 1)

		ck := http.Cookie{
			Name:     "auth_token",
			Domain:   "localhost",
			Path:     "/",
			Expires:  expires,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}

		// value of cookie
		ck.Value = token

		// write the cookie to response
		http.SetCookie(w, &ck)

		jsonBytesResponse(w, http.StatusOK, json)
		return
	}

	a.errorResponse(w, r.URL.Path, http.StatusBadRequest, "invalid login type", nil)
}

func (a *API) handleRegister(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /register register
	//
	// Register new user
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   description: Register request
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RegisterRequest"
	// responses:
	//   '200':
	//     description: success
	//   '401':
	//     description: invalid registration token
	//   '500':
	//     description: internal error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "failed to read request body", err)
		return
	}

	var registerData RegisterRequest
	err = json.Unmarshal(requestBody, &registerData)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "", err)
		return
	}
	registerData.Email = strings.TrimSpace(registerData.Email)
	registerData.Username = strings.TrimSpace(registerData.Username)

	err = a.app.RegisterUser(registerData.Username, registerData.Email, registerData.Password)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusBadRequest, err.Error(), err)
		return
	}

	jsonStringResponse(w, http.StatusOK, "{}")
}
