package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mertSezerO/HTTP-Server/database"
)

func (apiCfg apiConfig) endpointUsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		apiCfg.handlerUpdateUser(w, r)
	case http.MethodGet:
		apiCfg.handlerGetUser(w, r)
	case http.MethodPost:
		apiCfg.handlerCreateUser(w, r)
	case http.MethodDelete:
		apiCfg.handlerDeleteUser(w, r)
	default:
		respondWithError(w, 404, errors.New("method not supported"))
	}
}

func (apiCfg apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := database.User{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	user, err = apiCfg.client.CreateUser(user.Email, user.Password, user.Name, user.Age)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, 201, user)
}

func (apiCfg apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	email := strings.TrimPrefix(path, "/users/")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("no userEmail provided to delete"))
		return
	}
	err := apiCfg.client.DeleteUser(email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, 200, struct{}{})
}

func (apiCfg apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	email := strings.TrimPrefix(path, "/users/")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("no userEmail provided to get"))
		return
	}
	user, err := apiCfg.client.GetUser(email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, 200, user)
}

func (apiCfg apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	email := strings.TrimPrefix(path, "/users/")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("no userEmail provided to update"))
		return
	}
	decoder := json.NewDecoder(r.Body)
	user := database.User{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	updated, err := apiCfg.client.UpdateUser(email, user.Password, user.Name, user.Age)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, 200, updated)
}
