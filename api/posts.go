package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"solid-server/model"
	"solid-server/services/types"
)

func (a *API) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	// swagger:operation Post /posts/ create post
	//
	// Post Create
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   description: Post request
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
	ctx := r.Context()
	user := ctx.Value(sessionContextKey).(*model.User)

	_, err := json.Marshal(user)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "", err)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "failed to read request body", err)
		return
	}

	var registerData types.CreatePostRequest
	err = json.Unmarshal(requestBody, &registerData)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusInternalServerError, "", err)
		return
	}

	err = a.app.CreatePost(registerData, user.ID)
	if err != nil {
		a.errorResponse(w, r.URL.Path, http.StatusBadRequest, err.Error(), err)
		return
	}

	jsonStringResponse(w, http.StatusOK, "{}")
}
