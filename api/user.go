package api

import (
	"encoding/json"
	"net/http"
	"solid-server/model"
)

func (a *API) handleUserMe(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /users/me getMe
	//
	// Returns the currently logged-in user
	//
	// ---
	// produces:
	// - application/json
	// security:
	// - BearerAuth: []
	// responses:
	//   '200':
	//     description: success
	//     schema:
	//       "$ref": "#/definitions/User"
	//   default:
	//     description: internal error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	ctx := r.Context()
	user := ctx.Value(sessionContextKey).(*model.User)

	userData, err := json.Marshal(user)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "", err)
		return
	}

	jsonBytesResponse(w, http.StatusOK, userData)
}
