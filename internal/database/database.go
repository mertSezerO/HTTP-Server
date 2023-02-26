package database

import (
	"encoding/json"
	"errors"
	"os"
)

type Client struct {
	databasePath string
}
type databaseSchema struct {
	Users map[string]User `json:"users"`
	Posts map[string]Post `json:"posts"`
}

func NewClient(databasePath string) Client {
	return Client{
		databasePath: databasePath,
	}
}

func (c Client) createDatabase() error {
	dat, err := json.Marshal(databaseSchema{
		Users: make(map[string]User),
		Posts: make(map[string]Post),
	})
	if err != nil {
		return err
	}
	err = os.WriteFile(c.databasePath, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) ensureDatabase() error {
	_, err := os.ReadFile(c.databasePath)
	if errors.Is(err, os.ErrNotExist) {
		return c.createDatabase()
	}
	return err
}

func (c Client) updateDatabase(db databaseSchema) error {
	data, err := json.Marshal(db)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.databasePath, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (c Client) readDatabase() (databaseSchema, error) {
	data, err := os.ReadFile(c.databasePath)
	if err != nil {
		return databaseSchema{}, err
	}
	db := databaseSchema{}
	err = json.Unmarshal(data, &db)
	return db, err
}
