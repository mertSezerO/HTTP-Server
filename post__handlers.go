package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mertSezerO/HTTP-Server/database"
)

func (apiCfg apiConfig) endpointPostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		apiCfg.handlerRetrievePosts(w, r)
	case http.MethodPost:
		apiCfg.handlerCreatePost(w, r)
	case http.MethodDelete:
		apiCfg.handlerDeletePost(w, r)
	}
}

func (apiCfg apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	post := database.Post{}
	err := decoder.Decode(&post)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	post, err = apiCfg.client.CreatePost(post.UserEmail, post.Text)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, 201, post)
}

func (apiCfg apiConfig) handlerDeletePost(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/posts/")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("no id provided to get"))
		return
	}
	err := apiCfg.client.DeletePost(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, 200, struct{}{})
}

func (apiCfg apiConfig) handlerRetrievePosts(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("userEmail")
	if userEmail == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("no userEmail provided to retrieve posts of the account"))
		return
	}
	posts, err := apiCfg.client.GetPosts(userEmail)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, 200, posts)
}
