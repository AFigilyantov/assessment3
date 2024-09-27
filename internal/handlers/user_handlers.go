package handlers

import (
	"chitests/internal/models"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// trasport layer
// parse headers/body
// validation
// app call. make business logic
// render response

// TODO tranfer to another location i.e. folder http
type UserReposytory interface {
	RegisterUser(ctx context.Context, userAccount models.UserAccount) (models.UserAccount, error)
	LoginUser(ctx context.Context, username, password string) (models.UserAccount, error)
}

func RegisterHandler(u UserReposytory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req := &models.UserAccount{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {

			http.Error(w, "parsing err", http.StatusBadRequest)
			return
		}

		if err := validator.New().Struct(req); err != nil {
			http.Error(w, "validation err", http.StatusBadRequest)
			return
		}

		userAccount, err := u.RegisterUser(r.Context(), *req)

		if err != nil {
			http.Error(w, "registration faled", http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, userAccount)

	}
}

func LoginHandler(u UserReposytory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		req := &models.UserAccount{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, "parsing err", http.StatusBadRequest)
			return
		}

		userAccount, err := u.LoginUser(r.Context(), req.UserName, req.Password)

		if err != nil {
			http.Error(w, "error in handler", http.StatusNotFound)
			return
		}

		render.JSON(w, r, userAccount)

	}
}
