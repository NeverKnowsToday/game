package db

import (
	"github.com/game/server/database"
	"time"
)

// Token table
type Token struct {
	Token     string    `json:"token" gorm:"primary_key;column:token" form:"token"`
	Name      string    `json:"name"`
	Origin    string    `json:"origin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Create a new token
func InsertToken(name, token, origin string) error {
	return database.InsertDb(Db, &Token{Name: name, Token: token, Origin: origin})
}

// Get a token by token
func GetToken(token string) (*Token, error) {
	var t Token
	err := database.GetOne(Db, &t, "\"token\" = ?", token)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Delete a token
func DeleteToken(token string) error {
	return database.DeleteOne(Db, &Token{}, "\"token\" = ?", token)
}
