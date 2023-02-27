package database

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserEmail string    `json:"userEmail"`
	Text      string    `json:"text"`
}

func (c Client) CreatePost(userEmail, text string) (Post, error) {
	db, err := c.readDatabase()
	if err != nil {
		return Post{}, err
	}
	if _, exist := db.Users[userEmail]; !exist {
		return Post{}, errors.New("the user does not exists")
	}
	post := Post{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC(),
		UserEmail: userEmail,
		Text:      text,
	}
	db.Posts[post.ID] = post
	err = c.updateDatabase(db)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func (c Client) GetPosts(userEmail string) ([]Post, error) {
	db, err := c.readDatabase()
	if err != nil {
		return []Post{}, err
	}
	if _, exist := db.Users[userEmail]; !exist {
		return []Post{}, errors.New("the user does not exists")
	}
	posts := []Post{}
	for _, post := range db.Posts {
		if post.UserEmail == userEmail {
			posts = append(posts, post)
		}
	}
	return posts, nil
}

func (c Client) DeletePost(id string) error {
	db, err := c.readDatabase()
	if err != nil {
		return err
	}
	if _, exist := db.Posts[id]; !exist {
		return errors.New("the post does not exists")
	}
	delete(db.Posts, id)
	err = c.updateDatabase(db)
	if err != nil {
		return err
	}
	return nil
}
