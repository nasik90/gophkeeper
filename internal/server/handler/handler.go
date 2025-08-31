package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	middleware "github.com/nasik90/gophkeeper/internal/server/middlewares"
	"github.com/nasik90/gophkeeper/internal/server/storage"
)

type Service interface {
	RegisterNewUser(ctx context.Context, user, password string) error
	UserIsValid(ctx context.Context, login, password string) (bool, error)
	LoadSecret(ctx context.Context, key, value, login string) (int, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterNewUser() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var input struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.service.RegisterNewUser(ctx, input.Login, input.Password); err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, storage.ErrUserNotUnique) {
				status = http.StatusConflict
			}
			http.Error(res, err.Error(), status)
			return
		}
		if err := middleware.SetAuthCookie(input.Login, res); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) LoginUser() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var input struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		isValid, err := h.service.UserIsValid(ctx, input.Login, input.Password)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if !isValid {
			http.Error(res, "", http.StatusUnauthorized)
			return
		}
		if err := middleware.SetAuthCookie(input.Login, res); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) LoadSecret() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		login := middleware.LoginFromContext(ctx)
		var input struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		recordID, err := h.service.LoadSecret(ctx, input.Key, input.Value, login)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		var output struct {
			RecordID int `json:"recordID"`
		}
		output.RecordID = recordID
		outputJSON, err := json.Marshal(output)

		resStatus := http.StatusOK
		res.Header().Set("content-type", "application/json")
		res.WriteHeader(resStatus)
		res.Write(outputJSON)
	}
}
