package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"solid-server/services/auth"
	"strings"
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

// LoginRequest is a login request
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

// RegisterRequest is a user registration request
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

	// Registration authorization token
	// required: true
	Token string `json:"token"`
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
	if len(a.singleUserToken) > 0 {
		// Not permitted in single-user mode
		a.errorResponse(w, r.URL.Path, http.StatusUnauthorized, "not permitted in single-user mode", nil)
		return
	}

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

	//auditRec := a.makeAuditRecord(r, "login", audit.Fail)
	//defer a.audit.LogRecord(audit.LevelAuth, auditRec)
	//auditRec.AddMeta("username", loginData.Username)
	//auditRec.AddMeta("type", loginData.Type)
	return
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
	if len(a.singleUserToken) > 0 {
		a.errorResponse(w, r.URL.Path, http.StatusUnauthorized, "not permitted in single-user mode", nil)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "failed to read request body", err)
		return
	}

	var registerData RegisterRequest
	err = json.Unmarshal(requestBody, &registerData)
	err = json.Unmarshal(requestBody, &registerData)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "", err)
		return
	}
	registerData.Email = strings.TrimSpace(registerData.Email)
	registerData.Username = strings.TrimSpace(registerData.Username)

	// Validate Token
	if len(registerData.Token) > 0 {
		// TODO: Token Regiser
	} else {
		// 해당 토큰이 존재하는 경우 해당 토큰으로 가입한 유저가 있는지 체크
	}

	//auditRec := a.makeAuditRecord(r, "register", audit.Fail)
	//defer a.audit.LogRecord(audit.LevelAuth, auditRec)
	//auditRec.AddMeta("username", registerData.Username)

	err = a.app.RegisterUser(registerData.Username, registerData.Email, registerData.Password)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusBadRequest, err.Error(), err)
		return
	}

	jsonStringResponse(w, http.StatusOK, "{}")
	//auditRec.Success()
}
