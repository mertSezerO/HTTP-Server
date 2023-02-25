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

func createDatabase(c Client) error {
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

func ensureDatabase(c Client) error {
	_, err := os.ReadFile(c.databasePath)
	if errors.Is(err, os.ErrNotExist) {
		return createDatabase(c)
	}
	return err
}
