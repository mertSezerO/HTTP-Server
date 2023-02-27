package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mertSezerO/HTTP-Server/database"
)

func main() {
	m := http.NewServeMux()
	client := database.NewClient("database.json")
	apiCfg := apiConfig{
		client: client,
	}

	m.HandleFunc("/", testHandler)
	m.HandleFunc("/err", testErrHandler)
	m.HandleFunc("/users", apiCfg.endpointUsersHandler)
	m.HandleFunc("/users/", apiCfg.endpointUsersHandler)

	const addr = "localhost:8080"
	srv := http.Server{
		Handler:      m,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	// this blocks forever, until the server
	// has an unrecoverable error
	fmt.Println("server started on ", addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}

type apiConfig struct {
	client database.Client
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, database.User{
		Email: "test@example.com",
		Age:   25,
		Name:  "Test",
	})
}

func testErrHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, errors.New("server error"))
}

type errorLog struct {
	Error string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.WriteHeader(code)
	if payload != nil {
		response, err := json.Marshal(payload)
		if err != nil {
			log.Println("Error in Marshalling", err)
			w.WriteHeader(500)
			response, _ := json.Marshal(errorLog{
				Error: "error marshalling",
			})
			w.Write(response)
			return
		}
		w.WriteHeader(code)
		w.Write(response)
	}
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	if err == nil {
		log.Println("Error response with a nil error")
		return
	}
	log.Println(err)
	respondWithJSON(w, code, errorLog{
		Error: err.Error(),
	})
}

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
