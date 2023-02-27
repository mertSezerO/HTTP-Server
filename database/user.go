package database

import (
	"errors"
	"time"
)

type User struct {
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
}

func (c Client) CreateUser(email, password, name string, age int) (User, error) {
	db, err := c.readDatabase()
	if err != nil {
		return User{}, err
	}
	if _, exist := db.Users[email]; exist {
		return User{}, errors.New("the user already exists")
	}
	user := User{
		CreatedAt: time.Now().UTC(),
		Email:     email,
		Password:  password,
		Name:      name,
		Age:       age,
	}
	db.Users[email] = user
	err = c.updateDatabase(db)
	if err != nil {
		return User{}, err
	}
	return user, err
}

func (c Client) GetUser(email string) (User, error) {
	db, err := c.readDatabase()
	if err != nil {
		return User{}, err
	}
	if _, exist := db.Users[email]; !exist {
		return User{}, errors.New("the user does not exist")
	}
	return db.Users[email], nil
}

func (c Client) UpdateUser(email, password, name string, age int) (User, error) {
	db, err := c.readDatabase()
	if err != nil {
		return User{}, err
	}
	if _, exist := db.Users[email]; !exist {
		return User{}, errors.New("the user does not exist")
	}
	user := db.Users[email]
	user.Password = password
	user.Name = name
	user.Age = age
	err = c.updateDatabase(db)
	if err != nil {
		return User{}, err
	}
	return user, err
}

func (c Client) DeleteUser(email string) error {
	db, err := c.readDatabase()
	if err != nil {
		return err
	}
	if _, exist := db.Users[email]; !exist {
		return errors.New("the user does not exist")
	}
	delete(db.Users, email)
	err = c.updateDatabase(db)
	if err != nil {
		return err
	}
	return nil
}
